#import "@preview/bean-upm:0.1.0": upm-report
#import "@preview/fletcher:0.5.8" as fletcher: diagram, edge, node
#import "@preview/codly:1.3.0": *
#import "@preview/codly-languages:0.1.1": *
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


// Acknowledgements
#let acknowledgements = ""

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

  abstract-en: abst,

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

