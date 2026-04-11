-- +goose Up
-- +goose StatementBegin
-- =============================
-- EXTENSIONS
-- =============================
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS "citext";

-- =============================
-- ENUM TYPES
-- =============================
CREATE TYPE user_role_type AS ENUM ('STUDENT', 'INSTRUCTOR', 'ADMIN');
CREATE TYPE exam_status_type AS ENUM ('DRAFT', 'PUBLISHED', 'CLOSED');
CREATE TYPE subject_status_type AS ENUM ('ACTIVE', 'INACTIVE');
CREATE TYPE attempt_status_type AS ENUM ('IN_PROGRESS', 'SUBMITTED', 'GRADED', 'CANCELLED');
CREATE TYPE question_type_type AS ENUM ('TEXT', 'CODE', 'IMAGE');

-- =============================
-- FUNCTIONS
-- =============================
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION sync_subject_total_marks()
RETURNS TRIGGER AS $$
DECLARE
    target_subject_id UUID;
BEGIN
    IF TG_OP = 'DELETE' THEN
        target_subject_id := OLD.subject_id;
    ELSE
        target_subject_id := NEW.subject_id;
    END IF;

    UPDATE subjects
    SET total_marks = COALESCE((
        SELECT SUM(total_marks)
        FROM exams
        WHERE subject_id = target_subject_id
          AND deleted_at IS NULL
    ), 0)
    WHERE id = target_subject_id;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- =============================
-- USERS
-- =============================
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username CITEXT NOT NULL UNIQUE,
    full_name VARCHAR(255) NOT NULL,
    password_hash TEXT NOT NULL,
    role user_role_type NOT NULL DEFAULT 'STUDENT',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    last_login TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL,

    CHECK (char_length(username) >= 3),
    CHECK (username ~ '^[a-zA-Z0-9_]+$')
);

CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_active ON users(is_active);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

CREATE TRIGGER users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

-- =============================
-- SUBJECTS
-- =============================
CREATE TABLE subjects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    duration_minutes INT NOT NULL CHECK (duration_minutes > 0),
    total_marks INT NOT NULL DEFAULT 0 CHECK (total_marks >= 0),
    pass_score INT NOT NULL CHECK (pass_score >= 0),
    status subject_status_type NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL
);

CREATE TRIGGER subjects_updated_at
BEFORE UPDATE ON subjects
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

-- =============================
-- SUBJECT INSTRUCTORS (Join Table)
-- =============================
CREATE TABLE subject_instructors ( 
    subject_id UUID REFERENCES subjects(id) ON DELETE CASCADE,
    instructor_id UUID REFERENCES users(id) ON DELETE CASCADE,
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    assigned_by UUID REFERENCES users(id) ON DELETE SET NULL,
    deleted_at TIMESTAMPTZ,
    PRIMARY KEY (subject_id, instructor_id)
);

CREATE INDEX idx_subject_instructors_subject ON subject_instructors(subject_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_subject_instructors_instructor ON subject_instructors(instructor_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_subject_instructors_deleted ON subject_instructors(deleted_at) WHERE deleted_at IS NOT NULL;

-- =============================
-- SUBJECT STUDENTS (Join Table)
-- =============================
CREATE TABLE subject_students ( 
    subject_id UUID REFERENCES subjects(id) ON DELETE CASCADE,
    student_id UUID REFERENCES users(id) ON DELETE CASCADE,
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    assigned_by UUID REFERENCES users(id) ON DELETE SET NULL,
    deleted_at TIMESTAMPTZ,
    PRIMARY KEY (subject_id, student_id)
);

CREATE INDEX idx_subject_students_subject ON subject_students(subject_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_subject_students_student ON subject_students(student_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_subject_students_deleted ON subject_students(deleted_at) WHERE deleted_at IS NOT NULL;

CREATE TRIGGER subject_students_updated_at
BEFORE UPDATE ON subject_students
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

-- =============================
-- EXAMS
-- =============================
CREATE TABLE exams (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    subject_id UUID REFERENCES subjects(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    total_marks INT NOT NULL CHECK (total_marks > 0),
    status exam_status_type NOT NULL DEFAULT 'DRAFT',
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL
);

CREATE INDEX idx_exams_subject ON exams(subject_id);
CREATE INDEX idx_exams_id ON exams(id);
CREATE INDEX idx_exams_status ON exams(status);

CREATE TRIGGER exams_updated_at
BEFORE UPDATE ON exams
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER sync_subject_total_marks
AFTER INSERT OR DELETE ON exams
FOR EACH ROW
EXECUTE FUNCTION sync_subject_total_marks();

CREATE TRIGGER sync_subject_total_marks_on_update
AFTER UPDATE OF total_marks, deleted_at, subject_id ON exams
FOR EACH ROW
EXECUTE FUNCTION sync_subject_total_marks();

CREATE OR REPLACE FUNCTION validate_exam_creator_role()
RETURNS TRIGGER AS $$
DECLARE
    creator_role user_role_type;
BEGIN
    SELECT role INTO creator_role 
    FROM users 
    WHERE id = NEW.created_by;

    IF creator_role NOT IN ('INSTRUCTOR', 'ADMIN') THEN
        RAISE EXCEPTION 'User % is a %, and is not authorized to create or manage exams.', 
            NEW.created_by, creator_role;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER enforce_exam_creator_role
BEFORE INSERT OR UPDATE OF created_by ON exams
FOR EACH ROW
EXECUTE FUNCTION validate_exam_creator_role();

-- =============================
-- ENROLLMENTS
-- =============================
CREATE TABLE enrollments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID REFERENCES users(id) ON DELETE CASCADE,
    exam_id UUID REFERENCES exams(id) ON DELETE CASCADE,
    enrolled_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (student_id, exam_id)
);

CREATE INDEX idx_enrollments_student ON enrollments(student_id);
CREATE INDEX idx_enrollments_exam ON enrollments(exam_id);

-- =============================
-- QUESTIONS
-- =============================
CREATE TABLE questions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    exam_id UUID REFERENCES exams(id) ON DELETE CASCADE,
    question_title TEXT NOT NULL,
    question_text TEXT NOT NULL,
    question_type question_type_type NOT NULL DEFAULT 'TEXT',
    question_image TEXT,
    checksum TEXT,
    marks INT NOT NULL DEFAULT 1 CHECK (marks > 0),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL
);

CREATE INDEX idx_questions_exam ON questions(exam_id);
CREATE INDEX idx_questions_id ON questions(id);

CREATE TRIGGER questions_updated_at
BEFORE UPDATE ON questions
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

-- =============================
-- CHOICES
-- =============================
CREATE TABLE choices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    question_id UUID REFERENCES questions(id) ON DELETE CASCADE,
    choice_text TEXT NOT NULL,
    is_correct BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL
);

CREATE INDEX idx_choices_question ON choices(question_id);

CREATE TRIGGER choices_updated_at
BEFORE UPDATE ON choices
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE OR REPLACE FUNCTION check_choices_count()
RETURNS TRIGGER AS $$
DECLARE
    choice_count INT;
BEGIN
    SELECT COUNT(*) INTO choice_count 
    FROM choices 
    WHERE question_id = NEW.question_id;
    
    IF choice_count >= 4 THEN
        RAISE EXCEPTION 'Maximum 4 choices allowed per question';
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER check_choices_count
BEFORE INSERT ON choices
FOR EACH ROW
EXECUTE FUNCTION check_choices_count();

CREATE UNIQUE INDEX one_correct_choice_per_question
ON choices (question_id)
WHERE is_correct = true;

-- =============================
-- EXAM ATTEMPTS
-- =============================
CREATE TABLE exam_attempts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    exam_id UUID REFERENCES exams(id) ON DELETE CASCADE,
    student_id UUID REFERENCES users(id) ON DELETE CASCADE,
    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    submitted_at TIMESTAMPTZ,
    score INT CHECK (score >= 0),
    status attempt_status_type NOT NULL DEFAULT 'IN_PROGRESS',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    CHECK (submitted_at IS NULL OR submitted_at >= started_at)
);

CREATE INDEX idx_attempts_exam ON exam_attempts(exam_id);
CREATE INDEX idx_attempts_student ON exam_attempts(student_id);
CREATE INDEX idx_attempts_status ON exam_attempts(status);

CREATE UNIQUE INDEX one_active_attempt_per_exam
ON exam_attempts (exam_id, student_id)
WHERE status = 'IN_PROGRESS';

-- =============================
-- STUDENT ANSWERS
-- =============================
CREATE TABLE student_answers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    attempt_id UUID REFERENCES exam_attempts(id) ON DELETE CASCADE,
    question_id UUID REFERENCES questions(id) ON DELETE CASCADE,
    selected_choice_id UUID REFERENCES choices(id),
    is_correct BOOLEAN,
    answered_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (attempt_id, question_id)
);

CREATE INDEX idx_answers_attempt ON student_answers(attempt_id);
CREATE INDEX idx_answers_question ON student_answers(question_id);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS student_answers CASCADE;
DROP TABLE IF EXISTS exam_attempts CASCADE;
DROP TABLE IF EXISTS choices CASCADE;
DROP TABLE IF EXISTS questions CASCADE;
DROP TABLE IF EXISTS enrollments CASCADE;
DROP TABLE IF EXISTS exams CASCADE;
DROP TABLE IF EXISTS subject_instructors CASCADE;
DROP TABLE IF EXISTS subject_students CASCADE;
DROP TABLE IF EXISTS subjects CASCADE;
DROP TABLE IF EXISTS users CASCADE;

DROP FUNCTION IF EXISTS validate_exam_creator_role() CASCADE;
DROP FUNCTION IF EXISTS check_choices_count() CASCADE;
DROP FUNCTION IF EXISTS sync_subject_total_marks() CASCADE;
DROP FUNCTION IF EXISTS set_updated_at() CASCADE;

DROP TYPE IF EXISTS attempt_status_type CASCADE;
DROP TYPE IF EXISTS exam_status_type CASCADE;
DROP TYPE IF EXISTS user_role_type CASCADE;
DROP TYPE IF EXISTS subject_status_type CASCADE;
DROP TYPE IF EXISTS question_type_type CASCADE;

DROP EXTENSION IF EXISTS "citext";
DROP EXTENSION IF EXISTS "pgcrypto";
-- +goose StatementEnd