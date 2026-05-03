#import "@preview/fletcher:0.5.8" as fletcher: diagram, edge, node
#import "@preview/tiago:0.1.0": *
= System Design and \ Implementation

== Research Methodology
The development phase of this project was carried out using the *_Agile Methodology_* software development methodology. The system is a _web-based, data-driven application_. So, a *_design-driven development_* methodology was used to implement the system. At the analysis phase, the system requirements were gathered (*_functional and non-functional requirements_*) and analyzed.
A *_vertical slice_* of the system was implemented first into layers :

+ *The Foundation Layer*: database schema, domain model, and repositories.

+ *The Application Layer*: application services and use cases.

+ *The Presentation Layer*: web handlers and templates.


Security was a major concept in the design phase of the system; *TLS*, *CSRF*, *JWT*, *HTTP Security Headers*, and *RBAC* were implemented to ensure the security of the system.
The system was built using the *_Hexagonal Architecture_* pattern, which is a _domain-driven architecture_ that focuses on the separation of concerns and _dependency decoupling_. The stack was used in the system to give the system a modern, fast, _well-documented_, _well-maintained_, and _easy-to-scale stack_. focusing on *_performance_* and *_type-safety_*.

#pagebreak()
== System Design
=== System Architecture Overview
This system, as mentioned before, follows a Hexagonal architecture pattern, which separates code into separate domains:
\
\
#align(center, figure(
  diagram(
    spacing: 2em,
    node-corner-radius: 4pt,
    edge-stroke: 0.6pt,
    mark-scale: 70%,

    {
      let tint(c) = (stroke: c + 1pt, fill: c.lighten(95%), inset: 12pt)

      node(
        (0, 0),
        [
          *Inbound Adapters (HTTP)* \
          Router #sym.bar Handlers #sym.bar Middleware #sym.bar Helpers
        ],
        ..tint(teal),
        name: <inbound>,
      )

      edge(<inbound>, <app>, "-|>")

      node(
        (0, 1),
        [
          *Application Layer* \
          UserApplication #sym.bar ExamApplication #sym.bar SubjectAppl.
        ],
        ..tint(blue),
        name: <app>,
      )

      edge(<app>, <outbound>, "-|>", label: text(0.8em)[(via Port Interfaces)], label-pos: 0.5, label-side: left)

      node(
        (0, 2),
        [
          *Outbound Adapters* \
          PostgreSQL Repository #sym.bar Logger #sym.bar Config #sym.bar Storage
        ],
        ..tint(orange),
        name: <outbound>,
      )

      edge(<outbound>, <infra>, "-|>")

      node(
        (0, 3),
        [
          *Infrastructure* \
          PostgreSQL Database #sym.bar File System #sym.bar TLS Certs
        ],
        ..tint(gray),
        name: <infra>,
      )
    },
  ),
  caption: [System Architecture Overview],
))
\
\
- *Inbound Adapters*: Inbound adapters are often called *_driving adapters_*, because they drive the application. It is the only layer that interacts with the outside world. This system uses *_HTTP_* as the inbound adapter. This layer is responsible for handling the *_HTTP_* requests and responses.

- *Application Layer*: The application layer is responsible for the business logic of the system, and *_transaction management_*.

- *Outbound Adapters*: Outbound adapters are often called *_driven adapters_*, because they drive the application. This layer is responsible for driving the database, file system, logs, and configuration.

=== Technology Selection and Justification

#table(
  columns: (0.5fr, 1fr, 1fr),

  [Component], [Technology], [Justification],
  [Backend],
  [Go 1.25.6],
  [High-performance compiled language, modern ecosystem, type-safe, concurrency, ideal for enterprise-level applications],

  [HTTP Framework], [Gin], [High-performance, middleware support, well documented and maintained],
  [websocket \ support], [Gorilla Websocket], [Well-documented and maintained WebSocket library for Go],
  [Database],
  [PostgreSQL 18.3],
  [Relational database, ACID compliance, scalability, database triggers, and advanced data types],

  [Containerization], [Docker], [Industry standard for containerization],
  [Frontend], [HTMX], [Lightweight interactivity framework with Ajax capabilities],
  [Database \ driver], [pgx/v5], [Better Go driver for PostgreSQL than `database/sql`],
  [Query \ Builder], [SQLC], [Generate type-safe Go code from a SQL query, and supports repository pattern],
  [Templating \ Engine],
  [Templ],
  [Type-safe, compiles to Go code, secure from XSS attacks, and uses a component-based architecture],

  [Frontend \ interactivity],
  [HTMX],
  [Lightweight interactivity framework with Ajax capabilities, no or minimal JavaScript, server-side rendering],

  [CSS \ framework], [Tailwind CSS], [Utility-first CSS framework, well-documented and maintained],
  [Database \ Migrations], [Goose], [Simple command-line tool for database migrations],
  [Logging], [Uber's Zap], [High-performance, structured, and flexible configuration],
  [Log rotation], [Lumberjack], [log file rotation, compression, and flexible configuration],
  [Config], [Viper], [Flexible configuration management, supports multiple formats like `yaml`],
  [Auth], [JWT], [Stateless authentication],
  [Password \ Hashing], [Bcrypt], [Powerful hashing algorithm, that is the industry standard for password hashing],
  [Session \ management], [Gin-contrib/sessions], [Session management for CSRF tokens],
  [CSRF \ Protection], [Gin-contrib/sessions], [CSRF protection middleware for HTTP requests],
  [TLS \ Encryption], [net], [TLS support for HTTP/2],
  [PDF \ Generation], [Maroto], [Well documented and flexible PDF generator],
  [CSV \ Generation], [encoding/csv], [CSV generation for data imports and exports],
  [LaTeX \ Rendering], [KaTeX], [Fast and reliable LaTeX rendering],
  [Code \ Highlighting], [Highlight.js], [High-performance, many languages and themes supported],
)

=== Database Design and Schema
The database is built using *_PostgreSQL_*, the main focus of the database is to be *_normalized and efficient_*. The database utilizes PostgreSQL plugins such as *_pgcrypto_* for *_password hashing_* and *_uuid_* generation, and *_citext_* for *_case-insensitive text_*. Also, the utilization of *_database triggers_* and *_constraints_* ensures the integrity of the database.

#figure(
  ```sql
  CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  username CITEXT NOT NULL UNIQUE,
  full_name VARCHAR(255) NOT NULL,
  branch VARCHAR(255) NOT NULL,
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
  ```,
  caption: [User table schema],
)

#let code = ```
# Complete Database Schema Example
users: {
 shape: sql_table
 id: uuid {constraint: primary_key}
 username: citext {constraint: unique}
 full_name: varchar
 password_hash: text
 role: user_role_type
 is_active: boolean
 last_login: timestamptz
 created_at: timestamptz
}

subjects: {
 shape: sql_table
 id: uuid {constraint: primary_key}
 title: varchar {constraint: unique}
 description: text
 duration_minutes: int
 total_marks: int
 pass_score: int
 status: subject_status_type
}

subject_instructors: {
 shape: sql_table
 subject_id: uuid {constraint: foreign_key}
 instructor_id: uuid {constraint: foreign_key}
 assigned_by: uuid {constraint: foreign_key}
 assigned_at: timestamptz
}

subject_students: {
 shape: sql_table
 subject_id: uuid {constraint: foreign_key}
 student_id: uuid {constraint: foreign_key}
 assigned_by: uuid {constraint: foreign_key}
 assigned_at: timestamptz
}

exams: {
 shape: sql_table
 id: uuid {constraint: primary_key}
 subject_id: uuid {constraint: foreign_key}
 title: varchar {constraint: unique}
 total_marks: int
 status: exam_status_type
 created_by: uuid {constraint: foreign_key}
 created_at: timestamptz
}

enrollments: {
 shape: sql_table
 id: uuid {constraint: primary_key}
 student_id: uuid {constraint: foreign_key}
 exam_id: uuid {constraint: foreign_key}
 enrolled_at: timestamptz
}

questions: {
 shape: sql_table
 id: uuid {constraint: primary_key}
 exam_id: uuid {constraint: foreign_key}
 question_title: text
 question_type: question_type_type
 marks: int
}

choices: {
 shape: sql_table
 id: uuid {constraint: primary_key}
 question_id: uuid {constraint: foreign_key}
 choice_text: text
 is_correct: boolean
}

exam_attempts: {
 shape: sql_table
 id: uuid {constraint: primary_key}
 exam_id: uuid {constraint: foreign_key}
 student_id: uuid {constraint: foreign_key}
 started_at: timestamptz
 status: attempt_status_type
 score: int
}

student_answers: {
 shape: sql_table
 id: uuid {constraint: primary_key}
 attempt_id: uuid {constraint: foreign_key}
 question_id: uuid {constraint: foreign_key}
 selected_choice_id: uuid {constraint: foreign_key}
 is_correct: boolean
}

users.id <-> subject_instructors.instructor_id
subjects.id <-> subject_instructors.subject_id
users.id <-> subject_instructors.assigned_by

users.id <-> subject_students.student_id
subjects.id <-> subject_students.subject_id
users.id <-> subject_students.assigned_by

subjects.id <-> exams.subject_id
users.id <-> exams.created_by

users.id <-> enrollments.student_id
exams.id <-> enrollments.exam_id

exams.id <-> questions.exam_id
questions.id <-> choices.question_id

exams.id <-> exam_attempts.exam_id
users.id <-> exam_attempts.student_id

exam_attempts.id <-> student_answers.attempt_id
questions.id <-> student_answers.question_id
choices.id <-> student_answers.selected_choice_id
```.text

#align(center)[
  #figure(
    render(code),
    caption: [Database Entity-Relationship Schema],
  )
]

- *_Users_*: The users table is the central table of this system, which stores information about users that is used in authentication and authorization. The users can be a `STUDENT`, `INSTRUCTOR`, or `ADMIN`.

- *_Collections_*: The collections table is used to group exams together. The collections can be `ACTIVE` or `INACTIVE`.

- *_Exams_*: The exams table is responsible for storing information about exams. The exams can be in `DRAFT`, `PUBLISHED`, or `CLOSED` status. It will only be accessible to students if it is in `PUBLISHED` status and assigned to a collection that is `ACTIVE`.

- *_Questions_*: The questions table is used to store questions that are linked to a specific exam. That can be of `TEXT`, `CODE`, or `IMAGE` question type.

- *_Choices_*: The choices table is used to store the choices for the questions.

- *_Collection_instructors_*: The collection_instructors table is used to store a relationship between instructors and collections, which means the administrator can assign a collection to an instructor.

- *_Collection_students_*: The collection_students table is used to store a relationship between collections and students, similar to the collection_instructors table. The administrator can assign a collection to a student.

- *_Student_answers_*: The student_answers table is used to store the answers of students for each question at each exam attempt.

- *_Exam_attempts_*: The exam_attempts table is used to store information about exam attempts of students. It has three statuses: `IN_PROGRESS`, `SUBMITTED`, and `CANCELLED`.

=== User Roles and Access Control Mechanism
The system has three user roles: `STUDENT`, `INSTRUCTOR`, and `ADMIN`. Each role has its own set of privileges and user interface. These privileges can be listed as follows:

#table(
  columns: (0.4fr, 0.4fr, 1fr),
  [Role], [Interface], [Privileges],
  [Student],
  [Student \ Dashboard],
  [
    - View assigned exams
    - Attempt exams
    - View results
    - Export results as PDF
  ],

  [Instructor],
  [Instructor \ Dashboard],
  [
    - Create/edit/delete exams
    - add/edit/delete questions
    - assign/unassign collections to students
  ],

  [Admin],
  [Admin \ Dashboard],
  [
    - Create/edit/delete collections
    - add/delete/restore users
    - assign/unassign collections to instructors and students
    - Create/edit/delete exams
    - add/edit/delete questions
  ],
)

The mechanism by which the system checks if the user is authorized to access a specific resource or not is to use _HTTP middlewares_. In this system, the authorization is done by chaining two middleware functions:

+ `AuthenticationMiddleware()`: This middleware is responsible for checking the JWT token of the request and storing the claims in the gin context.

+ `RoleAuthMiddleware(allowedRoles...)`: This middleware takes the user role from the stored claims in the gin context and checks if the user is authorized to access the requested resource and if the user is active or not.

#figure(
  ```go
  func (m *AuthMiddleware) RoleAuthMiddleware(allowedRoles ...domain.UserRole) gin.HandlerFunc {
  roleSet := make(map[domain.UserRole]struct{}, len(allowedRoles))
  for _, r := range allowedRoles {
  roleSet[r] = struct{}{}
  }

  return func(c *gin.Context) {
  raw, exists := c.Get("claims")
  if !exists {
  c.Abort()
  return
  }
  if !claims.IsActive {
  c.Abort()
  return
  }
  userRole := domain.UserRole(claims.Role)
  if _, allowed := roleSet[userRole]; !allowed {
  c.Abort()
  return
  }
  c.Next()
  }
  }
  ```,
  caption: [Role-based authentication middleware],
)

==== System Flowchart

#align(center)[
  #figure(
    image("../assets/sys-flow.png"),
    caption: [System Workflow and User Interactions],
  )
]
#pagebreak()
=== System Configuration
To make the system modifiable and configurable, a config file system is used. The configuration is stored in a *YAML* file located at `/config/config.yaml`. The configuration has default values for all of the parameters if none of the parameters are provided to prevent any errors during the startup of the system.

#figure(
  ```yaml
  http_server:
   address: "localhost"
   port: "8443"
   tls_cert_dir: "certs"
   cert_file: "certs/cert.pem"
   key_file: "certs/key.pem"

  database:
   host: "localhost"
   port: 5432
   username: "testuser"
   password: "123456"
   database_name: "testdb"
   ssl_mode: "disable"
   max_connections: 25
   min_connections: 5
   max_conn_lifetime: 1h
   max_conn_idle_time: 30m

   // rest of the config file
  ```,
  caption: [System Configuration File],
)



#pagebreak()
=== Deployment Strategy
To deploy the system with a _modern, efficient, and scalable solution_, a containerization strategy is a must. The use of *docker* allows us to package and deploy the system in a dependency-free, isolated environment.

#let docker_code = ```
docker: Docker Environment {
 network: Internal Network {
 app: Go Backend Application
 db: PostgreSQL Database
 }
}

client: Internet Client

client -> docker.network.app
docker.network.app <-> docker.network.db
```.text

#align(center)[
  #block(width: 300pt)[
    #figure(
      render(docker_code),
      caption: [Docker Container Deployment Architecture Workflow],
    )
  ]
]



== System Implementation
To provide a comprehensive understanding of the implementation of the system, we can trace a single event through the system, such as getting a list of all users in the system by utilizing a method called `GetUsers()`, for example.

=== The Domain Layer
In the domain layer located at `/internal/domain`, we define the entities and business logic of the system with no dependencies or knowledge of the outside world.


#align(center)[
  #block(width: 300pt)[
    #figure(
      ```go
      type User struct {
      ID        string    `json:"id"`
      Username  string    `json:"username"`
      FullName  string    `json:"full_name"`
      Role      UserRole  `json:"role"`
      IsActive  bool      `json:"is_active"`
      CreatedAt time.Time `json:"created_at"`
      UpdatedAt time.Time `json:"updated_at"`
      }
      ```,
      caption: [Domain Layer Entity],
    )
  ]
]

=== The Port Layer
The Port layer `/internal/ports` works as an interface between the domain layer and the outside world. That way the domain layer is not dependent or aware of the outside world.

#align(center)[
  #block(width: 300pt)[
    #figure(
      ```go
      type UserRepository interface {
      GetUsers() ([]User, error)
      }
      ```,
      caption: [Port Layer Interface],
    )
  ]
]

=== The Infrastructure Layer (Outbound Adapter)
This layer `/internal/adapters/outbound` is the lowest level of implementation. This layer contains the implementation of the ports defined in the port interface. In this system, *SQLC* is already generating the repositories for us, but the system is now dependent on *SQLC* to solve this dependency issue. A better solution is to map the generated sqlc structs to the domain structs in another repository layer.

```go
func (r *UserRepository) ListAll(ctx context.Context, arg ports.ListAllUsersParams) ([]domain.User, error) {
 queries := r.queries
 if tx := database.ExtractTx(ctx); tx != nil {
 queries = queries.WithTx(tx)
 }

 users, err := queries.ListAllUsers(ctx, sqlc.ListAllUsersParams{
 Limit:  arg.Limit,
 Offset: arg.Offset,
 })
 if err != nil {
 return nil, err
 }

 var domainUsers []domain.User
 for _, user := range users {
 domainUsers = append(domainUsers, domain.User{
 ID:        user.ID,
 Username:  user.Username,
 FullName:  user.FullName,
 Role:      domain.UserRole(user.Role),
 IsActive:  user.IsActive,
 LastLogin: toTimePtr(user.LastLogin),
 CreatedAt: user.CreatedAt.Time,
 UpdatedAt: user.UpdatedAt.Time,
 DeletedAt: toTimePtr(user.DeletedAt),
 })
 }

 return domainUsers, nil
}
```
This repository method satisfies the port interface.

=== The Application Layer
The application layer `/internal/adapters/inbound` is the wrapper around the port interface and the highest level of implementation of the system. The application layer doesn't have any knowledge of what storage or what database is being used.

```go
func (app *Application) ListAllUsers(ctx context.Context, arg ports.ListAllUsersParams) ([]domain.User, error) {
 users, err := app.userRepo.ListAll(ctx, arg)
 if err != nil {
 return nil, err
 }
 return users, nil
}
```

=== The Presentation Layer (Inbound Adapter)
The presentation layer `/internal/adapters/inbound` is the layer that directly interacts with the outside world. It deals with the incoming and outgoing data of the system.

```go
func (h *UserHandler) ListAllUsers() gin.HandlerFunc {
 return func(ctx *gin.Context) {
 users, err := h.App.ListAllUsers(ctx, ports.ListAllUsersParams{Limit: 100, Offset: 0})
 if err != nil {
 helpers.Toast(ctx, "List All Users Failed", "Failed to list all users: " + err.Error(), toast.VariantError)
 return
 }

 render.Render(ctx, components.UserTableRows(users, false))
 }
}
```

== User Interface
This system is built using a component-based architecture with *Templ* and *HTMX*. These components can be triggered by the handlers in the _presentation layer_. For example, a *_Toast Notification_* is a component that can be triggered from the handler using *HTMX AJAX* Headers `HX-Trigger` and `HX-Reswap`.

#figure(
  ```go
  func (h *UserHandler) ListAllUsers() gin.HandlerFunc {
  return func(ctx *gin.Context) {
  users, err := h.App.ListAllUsers(ctx, ports.ListAllUsersParams{Limit: 100, Offset: 0})
  if err != nil {
  helpers.Toast(ctx, "List All Users Failed", "Failed to list all users: " + err.Error(), toast.VariantError)
  return
  }
  render.Render(ctx, components.UserTableRows(users, false))
  }
  }
  ```,
  caption: "Toast Notification In Presentation Layer",
)

#figure(
  ```go
   func Toast(ctx *gin.Context, title string, description string, variant toast.Variant) {
   ctx.Writer.Write([]byte(`<div hx-swap-oob="beforeend:#toast-container">`))

   toast.Toast(toast.Props{
   Title:         title,
   Description:   description,
   Variant:       variant,
   Duration:      4000,
   ShowIndicator: true,
   Dismissible:   true,
   Icon:          true,
   }).Render(ctx.Request.Context(), ctx.Writer)

   ctx.Writer.Write([]byte(`</div>`))
  }
  ```,
  caption: "Toast Notification Implementation",
)

#figure(
  image("../assets/toast.png", width: 80%),
  caption: "Toast Notification After Triggering",
)


== Security Mechanisms
=== OWASP Top 10 Security Considerations
OWASP Top 10 is a list of the most important and critical security risks for web applications @owasp2025. The system is designed with these risks mitigation in mind.
\
\
#table(
  columns: (1fr, 1fr),
  [OWASP Top 10 Risk], [Mitigation Strategy],
  [A01:2025 Broken Access Control],
  [Role-Based Access Control (RBAC) middleware verifying allowed roles from JWT claims.],

  [A02:2025 Security Misconfiguration],
  [Strict TLS encryption enforcement, HTTP security headers (`X-Content-Type-Options`, `Frame-Options`), and configuration via Viper.],

  [A03:2025 Software Supply Chain Failures],
  [Use of Go's built-in module management (`go.mod`, `go.sum`) to strictly lock and verify dependency checksums.],

  [A04:2025 Cryptographic Failures],
  [Bcrypt algorithm for resilient password hashing; robust JWT token signing; enforcing `HttpOnly` and `Secure` cookie flags over HTTPS.],

  [A05:2025 Injection],
  [SQLC type-safe prepared statements natively eliminating SQL injection; Templ strictly escaping all dynamic HTML inputs against XSS.],

  [A06:2025 Insecure Design],
  [Hexagonal architecture isolating business logic from external frameworks, fortified with strict layer boundary separation.],

  [A07:2025 Authentication Failures],
  [Stateless JWT tokens with rigorous expiration deadlines; automated CSRF protection using `gin-csrf` rotating tokens.],

  [A08:2025 Software or Data Integrity Failures],
  [Data integrity enforced via database `CHECK` constraints, unique keys, structural triggers, and soft-deletion tracking.],

  [A09:2025 Security Logging and Alerting \ Failures],
  [Persistent, leveled, and structured JSON audit logging utilizing Uber's Zap logger matched with Lumberjack log rotation.],

  [A10:2025 Mishandling of Exceptional \ Conditions],
  [Centralized robust error-wrapping across all Adapters to prevent raw DB driver stack traces from leaking to clients.],
)

=== SQL Injection Prevention Techniques
SQL injection is a critical vulnerability that exploits the escape characters in the database queries @hu2020survey. This system uses type-safe prepared statements generated by *SQLC* to prevent SQL injection attacks. The *SQLC* takes a `.sql` file and generates a type-safe Go code with a repository layer that used to interact with the database.

=== TLS Encryption
*HTTP* traffic goes through *HTTPS* with the help of *TLS* certificates located at cert.pem and key.pem in certs/. The implementation that *Go* standard library uses natively (*net/tls*), which is compatible with TLS 1.2 and TLS 1.3. The *HTTP* response header `Strict-Transport-Security max-age=31536000; includeSubDomains` is added to ensure that browsers will exclusively connect to the domain and its subdomains using *HTTPS* for a year from the first visit (*HSTS*).

The `Referrer-Policy: strict-origin-when-cross-origin` header limits the referrer information sent in cross-origin requests, preventing the leakage of sensitive URL parameters to external sites.

The system uses a self-signed certificate for local deployment and the ability to use a signed certificate from the configuration *YAML* file.

```go
 privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
 if err != nil {
 return "", "", fmt.Errorf("failed to generate private key: %w", err)
 }

 template := x509.Certificate{
 SerialNumber: big.NewInt(1),
 Subject: pkix.Name{
 Organization: []string{"University of Technology"},
 Country:      []string{"IQ"},
 Locality:     []string{"baghdad"},
 },
 NotBefore:             time.Now(),
 NotAfter:              time.Now().Add(365 * 24 * time.Hour),
 KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
 ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
 BasicConstraintsValid: true,
 }
```
=== System Logging and Monitoring
The system utilizes structured logging with `Uber's Zap` library to log all the system events. The logs are stored in a file with a rotation mechanism to prevent the log file from growing indefinitely using `lumberjack` package, which is a compression and rotation mechanism for log files. *Zap* is a _near-zero allocation, fast, structured_ logging library for *Go*.

#figure(
  ```json
   {
   "level":"INFO",
   "timestamp":"2026-03-08T23:31:25.732+0300",
   "caller":"server/server.go:74",
   "message":"INTERNAL_INFO",
   "info_msg":"Using an existing TLS certificate..."
   }
  ```,
  caption: [System Log Example],
)

The use of *JSON* format for logging makes it easier to parse and analyze the logs using tools like *ELK Stack* or *Splunk*.

The system also has a separate log for all the routes and the *HTTP* requests.

#figure(
  ```json
  [GIN] 2026/04/14 - 10:02:56 | 200 | 4.992395ms | 127.0.0.1 | GET "/admin/dashboard/users"
  [GIN] 2026/04/14 - 10:02:59 | 200 | 3.799645ms | 127.0.0.1 | GET "/admin/dashboard/"
  [GIN] 2026/04/14 - 10:03:01 | 200 | 2.919365ms | 127.0.0.1 | GET "/admin/dashboard/subjects"
  ```,
  caption: [System Route Log Example],
)
=== Password Security
The system uses the *bcrypt* algorithm for password hashing. *bcrypt* is a cryptographic password hashing function based on the _Blowfish cipher_, designed by _*Niels Provos*_ and _*David Mazières*_ to protect against brute-force attacks. It is considered the industry standard for password hashing. *bcrypt* can be configured with a _cost factor_, the higher the cost factor, the more secure the password is, but the longer it takes to hash the password.

#figure(
  image("../assets/bcrypt.png", width: 100%),
  caption: [Bcrypt Password Hashing Algorithm],
)
\
\
The system also validates the password strength using a set of rules that are defined in the system. The rules are as follows:

- Password must be at least 8 characters long.
- Password must contain at least one number.
- Password must contain at least one uppercase letter.

#figure(
  ```go
  func ValidatePasswordCriteria(password string) error {
   if password == "" {
   return fmt.Errorf("password is required")
   }
   if len(password) < 8 {
   return fmt.Errorf("password must be at least 8 characters long")
   }
   var hasNum, hasUpper bool
   for _, char := range password {
   if unicode.IsNumber(char) {
   hasNum = true
   }
   if unicode.IsUpper(char) {
   hasUpper = true
   }
   }
   if !hasNum || !hasUpper {
   return fmt.Errorf("password must contain at least one number and one uppercase letter")
   }
   return nil
  }
  ```,
  caption: [Password Validation Function],
)

