#import "lib.typ": upm-report
#import "@preview/fletcher:0.5.8" as fletcher: diagram, edge, node
#import "@preview/codly:1.3.0": *
#import "@preview/codly-languages:0.1.1": *
#import "@preview/auto-bidi:0.1.0": *
#show: auto-dir
#show: codly-init.with()
#codly(
  zebra-fill: none,
  fill: rgb("#f5f5f5"),
  radius: 4pt,
  inset: 3pt,
  languages: codly-languages,
)


// Abstract
#let abst = [
  The reliance on digital infrastructure in academic institutions is rapidly increasing. The need for a *robust, secure, scalable, and performant examination system* is inevitable. This thesis tries to represent a design and implementation to solve such a case as a _web-based examination management system application_, keeping in mind all the mentioned aspects, in addition to a clean and modern architecture, a clean and easy interface, easy to scale and deploy with modern tooling, flexibility of question symbols and formats, and a role-aware digital platform. aimed to replace the manual process of traditional examinations.

  The system supports three distinct roles (_Administrator_, _Instructor_, and _Student_), each with its own set of capabilities and privileges, utilizing role-based access control as an *_HTTP middleware_* and *database constraints*. The administrator is responsible for managing exams, subjects, user assignments, and deletions through a dedicated dashboard. As for the instructors, they are responsible for the creation of the question in the form of an *MCQ question* with three types of questions:
  *TEXT*: which can deal with mathematical notations using $"LaTeX"$ directly in the browser.
  *CODE*, which can highlight a code snippet according to the written programming language.
  *IMAGE*: This is the final type that can add an image to a question.
  As for the Student role, after the admin enrolls the student successfully, they can simply enter the exam or the exam collection and submit before the end of the time duration; otherwise, the attempt will be auto-submitted, and they will receive automatically graded scores upon completion that can be downloaded as a *PDF format*.

  This framework is primarily built in the *_Go programming language_* as its backend language, using a clean architecture pattern known as _Hexagonal Architecture (Port and Adapters)_, a _domain-driven architecture_ that focuses on the separation of concerns and dependency decoupling. As for the frontend it is rendered using a type-safe compiled go template engine called *_Templ_*, with a _lightweight_ interactivity framework with *Ajax (Asynchronous JavaScript And XML)* capabilities and minimal _javascript_ code. Data persistence is handled by *_PostgreSQL_*, with database integrity enforced through custom triggers, constraints, and stored procedures.

  Security measures embedded in the system include _*JWT-based* cookie authentication_, _role-based access control_, *CSRF* protection, *bcrypt* password hashing, *TLS encryption*, HTTP security headers, and comprehensive structured logging using *Uber's Zap* library. The architecture prioritizes correctness and data integrity at every layer, from the application service down to the relational schema.
]

#let keywords-en = "Go, PostgreSQL, Templ, HTMX, Clean Architecture, Web Security, Examination Management System"
#let keywords-ar = "جو، بوستجريس، تيمبل، اتش تي ام اكس، الهندسة النظيفة، أمن الويب، نظام إدارة الامتحانات"

#let title-ar = "تصميم وتنفيذ نظام امتحانات فعال و امن"
#let author-ar = "الحسن عدي صاحب ، نور حسن عبد المهدي"
#let supervisor-ar = "أ.د. ضاري عادل محمود"
#let university-ar = "الجامعة التكنولوجية"
#let college-ar = "كلية هندسة الحاسوب"
#let department-ar = "قسم هندسة الحاسوب"
#let date-ar = "2025-2026"

#let abstract-ar = text(size: 15pt)[
تزايد الاعتماد على البنية التحتية الرقمية في المؤسسات الأكاديمية بشكل متسارع، مما جعل الحاجة إلى نظام امتحانات قوي، وآمن، وقابل للتوسع أمراً حتمياً. يهدف هذا البحث إلى تقديم تصميم وتنفيذ لتطبيق ويب لإدارة الامتحانات، مع مراعاة كافة الجوانب التقنية الحديثة، بما في ذلك الهندسة البرمجية النظيفة، وواجهة الاستخدام السهلة، ومرونة تنسيق رموز الأسئلة. يسعى النظام ليكون بديلاً رقمياً متطوراً للعمليات اليدوية في الامتحانات التقليدية، مما يضمن كفاءة الأداء وسهولة النشر.

يدعم النظام ثلاثة أدوار رئيسية (المسؤول، المحاضر، والطالب)، حيث يتمتع كل دور بصلاحيات محددة عبر نظام تحكم بالوصول مبني على قواعد بيانات وقيود برمجية. يتولى المسؤول إدارة الامتحانات والمواد والمستخدمين عبر لوحة تحكم مخصصة، بينما يقوم المحاضرون بإنشاء أسئلة الاختيار من متعدد بأنواع مختلفة؛ تشمل النصوص التي تدعم الصيغ الرياضية (بإستخدام LaTeX)، والأكواد البرمجية، والصور. أما الطلاب، فيمكنهم الدخول للامتحانات وتسليم الإجابات قبل انتهاء الوقت، مع ميزة التسليم التلقائي عند انتهاء المدة، والحصول الفوري على نتائج مصححة آلياً يمكن تحميلها بصيغة PDF.

من الناحية التقنية، تم بناء الواجهة الخلفية للنظام بلغة البرمجة Go وفق نمط "الهندسة السداسية" (Hexagonal Architecture) لضمان فصل المهام واستقلالية التبعيات. تعتمد واجهة المستخدم على محرك قوالب Templ مع تقنيات Ajax لضمان التفاعل السريع، بينما تتم إدارة البيانات عبر PostgreSQL مع فرض نزاهة البيانات من خلال القيود والإجراءات المخزنة. كما يتضمن النظام تدابير أمنية شاملة تشمل التوثيق عبر JWT ،وتشفير كلمات المرور بآلية Bcrypt، وتشفير TLS، لضمان أعلى مستويات الأمان وحماية البيانات في كل طبقات التطبيق.

]





// Acknowledgements
#let acknowledgements = text(size: 15pt)[
  We would like to express our sincere gratitude and profound appreciation to our supervisor, *Prof. Dhari A. Mahmood*, for his invaluable guidance, continuous support, and insightful feedback throughout the development of this project. His expertise and encouragement were instrumental in shaping the direction of this research.

  Our thanks also extend to the *Computer Engineering Department* at the *University of Technology* for providing the academic environment and resources necessary for this study. We are grateful to all the professors and staff who have shared their knowledge and experience with us during our years of study.

  Finally, we owe a debt of gratitude to our families for their unconditional love, patience, and unwavering support. This achievement would not have been possible without their constant encouragement and belief in our potential.
]

#let supervisorCert = text(size: 15pt)[
  I certify that the preparation of this project entitled (#text(fill: blue)[Design And Implementation of an Efficient and Secure Examination Management Framework]) was prepared by (#text(fill: blue)[Al-hassan Oday Sahib & Noor Hassan Abdulmahdi]) under my supervision at the Computer Engineering Department, College of Computer Engineering, the University of Technology in partial fulfillment of the requirements for the degree of B.Sc. in Computer Engineering.
]

#let examCert = text(size: 15pt)[
  We certify that we have read this project entitled (#text(fill: blue)[Design And Implementation of an Efficient and Secure Examination Management Framework]) , and as an examining committee examined the students (#text(fill: blue)[Al-hassan Oday Sahib & Noor Hassan Abdulmahdi]) in its contents, and our opinion, it meets the standards of the B.Sc. Computer Engineering.
]

#show: upm-report.with(
  title: "Design And Implementation of an Efficient and Secure Examination Management Framework",
  author: "Al-hassan Oday Sahib , Noor Hassan Abdulmahdi",
  supervisor: "Prof. Dhari A. Mahmood",
  date: datetime(year: 2026, month: 4, day: 18),

  university: "University of Technology - Iraq",
  school-name: "College of computer engineering - Computer Networks Engineering Department",
  school-abbr: "UOT",
  report-type: "Bachelor's thesis",
  degree-name: "",

  acknowledgements: acknowledgements,

  supervisorCert: supervisorCert,

  abstract-en: abst,
  keywords-en: keywords-en,
  abstract-ar: abstract-ar,
  keywords-ar: keywords-ar,

  title-ar: title-ar,
  author-ar: author-ar,
  supervisor-ar: supervisor-ar,
  university-ar: university-ar,
  college-ar: college-ar,
  department-ar: department-ar,
  date-ar: date-ar,

  ExaminationCert: examCert,

  school-logo: image("assets/logoPNG-White.png", width: 55mm),
  school-watermark: move(image("assets/logoPNG-Watermark.png", width: 150mm), dy: 220pt, dx: -20pt),
)


// Chapter 1
#include "chapters/chapter1.typ"

// Chapter 2
#include "chapters/chapter2.typ"

// Chapter 3
#include "chapters/chapter3.typ"

// Chapter 4
#include "chapters/chapter4.typ"

// Chapter 5
#include "chapters/chapter5.typ"

// Appendix
#include "chapters/Appendix.typ"

// Bib
#bibliography("references.bib", style: "ieee")

