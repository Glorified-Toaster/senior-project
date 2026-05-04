#set text(font: "Times New Roman")
= Introduction

== Overview
This project is the design and implementation of a _full-stack web-based application_ for an academic examination system. This system is intended to overcome the limitations of the traditional pen and paper examination approach. Also, it is intended to provide solutions for the problems of the existing digital solutions. The system is implemented using a _clean architecture approach_ for a scalable and maintainable solution. It uses *Go* on the backend, *PostgreSQL* on the database, and *Templ* and *HTMX* on the frontend. It incorporates a _layered security model_ that implements *JWT-based authentication*, *Role-Based Access Control (RBAC)*, *database constraints*, *HTTP Security Headers*, *CSRF tokens*, and *TLS encryption*. Also, support for *$"LaTeX"$* math formulas and code blocks. The system is evaluated with a set of benchmarks like performance, resource utilization, and security. Then, the results are compared with other existing solutions.

== Problem Statement
Academic Institutions continue to have challenges in conducting a secure, efficient examination, so the conventional way to do it is to heavily rely on the physical pen and paper approach. which comes with numerous operational problems, high paper and printing costs, physical logistics of distributing and collecting exam booklets, time-consuming manual grading, risk of loss or damage to submitted papers, and limited scalability as the student population grows @pariksha-2022.

Furthermore, the existing digital solution is lacking proper properties and features such as _role-based access control, tamper-proof answer recording, automated scoring, audit trails, performance, and modern tooling_ @ai-based-proctoring2025. Nevertheless, the lack of a good user interface and text formatting constraints.
#pagebreak()

== Aim of the Project
the primary aim of this project is to design and implement a performant, secure, scalable, maintainable, easy-to-use, and *_feature-rich web-based application_*, both in the _backend_ and _frontend_. this system is intended to serve as a replacement for the manual and slower academic examination workflows.

== Scope
+ Create and Deploy a full-stack web-based application to manage university exams that includes server, database schema, deployment, and server-side rendered frontend.

+ Use Hexagonal Architecture (Port and Adapters) to create a clean architecture design that separates domains and decouples dependencies through abstractions which allows for a testable, scalable, and maintainable system that is capable of easily adapting to future needs.

+ Create a layered security model that implements *JWT-based authentication*, *Role-Based Access Control (RBAC)*, *database constraints, HTTP Security Headers, CSRF tokens, and TLS encryption*.

+ Leverage *PostgreSQL*’s scripting language and automation through the use of stored functions and triggers to manage data consistency.

+ Reduce database latency through the use of _Database Indexing_.

== Objective
+ Design the database schema.

+ Implement the *RESTful API* and web server using Go, with a clean code hexagonal architecture and dependency injection pattern.

+ Implement role-based access control as a database trigger and HTTP middleware.

+ Implement logging and log rotation using Zap and Lumberjack.

+ Implement server best practices like graceful shutdown.

+ Implement the frontend with a server-side rendering interface using Templ and HTMX.

+ Implement $"LaTeX"$ support and syntax highlight.

+ Implement a security model that solves the *OWASP Top 10* vulnerabilities.

+ Evaluate the system performance and utilization and compare it with another traditional architectural approaches.

== Organization of the Project

- *Literature Review*: Provides a background, covering e-assessment systems, server-side rendering, database design patterns, architectural approaches, academic integrity mechanisms, and an analysis of related existing systems.

- *System Design and Implementation*: Describes the system design and implementation in detail, including the system architecture, technology selections, database schema, security mechanisms, and the implementation of core features.

- *Evaluation*: Presents an evaluation of the system, including architectural comparisons, performance measurements, and resource utilization analysis.

- *Conclusion and Future Work*: Concludes the thesis with a summary of contributions and directions for future work.

- *Appendices*: Provide source code snippets, the complete database schema, and user interface screenshots.
