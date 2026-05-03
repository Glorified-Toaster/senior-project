= Literature Review
== E-Assessment Systems
*Electronic Assessment* (E-Assessment) or *compouter-based assessment*, is the use of _digital technologies_ to measure and evaluate student learning outcome.
The *e-assessment systems* have been on the rise for the past few years, caused by the rapid development and innovation in the _computer network_ field, and the rise of the *_internet_*. Providing a replacement for the traditional _paper-based_ exams.

=== The advantages of such a platform
+ The reduction or the removal of the administrative overhead.

+ The reduction of the printing costs and paper prices.

+ The preservation of data that is very hard or impossible to tamper with.

+ Flexible and Scalable Solution.

+ Grading time is not affected by the number of students.

== Evolution of E-Assessment Systems
The first generation of computer-based exams was almost exclusively offline because of the lack of computer connectivity back then, to
Make the examination software up and running; a school or an institution must install the software locally with the aid of a specialist in the software. That means the software was limited and costly.

After the rise of the *WWW (World Wide Web)*, an online examination system began to take shape with simple security measures like *RBAC (Role-Based Access Control)* @edugauard-2025.

Over the years, advances in computer-based online examinations have been remarkable and fast, shifting focus from the simple examination platform to platforms that are trying to solve problems like _Security, Scalability, and automation_.

In the *COVID-19* pandemic, the need for a _secure, convenient, automated, and performant_ examination system was critical @ai-based-proctoring2025. examination management systems increasingly adopt _cloud-native architectures, microservice decomposition, modern frontend and backend frameworks to address scalability, maintainability, and user experience challenges_
== Static sites vs. Dynamic sites
=== Static sites
static site page or a flat page is a web page that is deliviered to the end user without any modification on the page file content.

#figure(
  image("../assets/static_app_server.png", width: 150mm),
  caption: "Static Sites Architecture",
)
=== Dynamic sites
dynamic sites are web pages that are generated according to the user request, it aggregates data from the database and renders it.

#figure(
  image("../assets/dynamic_app_server.png", width: 150mm),
  caption: "Dynamic Sites Architecture",
)

The use of dynamic sites is more sutiable for a data-driven application such as an electronic examination management system, it allows for more flexibility and interactivity.

#pagebreak()

== Type-Safe Templating in Web Applications
The *_template engine_* is the base of building a dynamic _server-side rendered_ web application. Template engines are crucial for separating _domain logic_ from _the presentation_. The output of the template engine is usually *HTML* or *XML* @Pisu_2026.

=== Traditional Templating Engines
The traditional template engine that is used in programming languages like *Python* (*Jinja*), *PHP* (*Twig*), *Ruby* (*ERB*) uses the string interpolation technique to generate content at _runtime_. This technique is prone to security vulnerabilities like *XSS* (*cross-site scripting*), *SSTI* (*Server-side template injection*). The fallowing code snippet is an example of a vulnerable *SSTI* template engine in *python*:

#figure(
  ```python
  from flask import Flask, request, render_template_string

  app = Flask(__name__)

  @app.route('/greet')
  def greet():
      name = request.args.get('name', 'Guest')
      template = f"<h1>Hello {name}</h1>"
      return render_template_string(template)
  ```,
  caption: "An example of an SSTI vulnerable Python template engine source code.",
)

It is quite easy to misuse the string interpolation technique. A caution is needed when developing with a type of template engine @Pisu_2026.
=== Type-Safe Templating Engines
On the contrary, the alternative to the traditional template engines is much stricter and more secure, by implication of _strong-typing_ and stricter escaping rules, such as *Rust*'s (*askama*) and *Go*'s (*Templ*), which are _compiled optimized template engines_.

#figure(
  ```go
  package main

  templ Hello(name string) {
    <div>Hello, { name }</div>
  }

  templ Greeting(person Person) {
    <div class="greeting">
      @Hello(person.Name)
    </div>
  }
  ```,
  caption: "An example of a type-safe template engine in Go (Templ)",
)
*Templ* is a _type-safe compiled_ template engine that compiles the template files into *Go code*, then compiles it to a *binary*. This pattern gives a huge advantage; these advantages @Hoenig2023Templ include:

+ *Compile-time* type-safety.

+ Automic *XSS* prevention.

+ Use of the *Go* _ecosestem and composition pattern_.

+ Performance and optimization.

== Server-Side Rendering vs. Single-Page Applications
In modern web development, there are two main approaches to developing a web application: the first approach is the traditional _server-side rendering_, and the second approach, and the most dominant one in modern web development, is *SPA (Single Page Application)* using _JavaScript frameworks_ like *Vue.js, React.js, and Angular.js*. While SPA provides a better _ecosystem_ and more interactivity due to the _dynamic compiled_ updating without reloading the full page, it throws the weight of computation to the client to handle it. *SPA* are also more complex to implement with more _JavaScript bundles_; this complexity can also produce a _security risk and critical vulnerabilities_ @cve_2025_55182 @cve_2025_66412.


While *SSR (Server-Side Rendering)* is a more traditional approach, it is still widely used in many applications, especially in _data-driven applications_ like examination management systems. SSR frameworks like *Ruby on Rails*, and *ASP.NET*, focusing on _performance, security, and ease of implementation_. A _lightweight framework_ like *HTMX* can be a middle ground between classic *SSR* and *SPA* with minimal JavaScript. *HTMX* can emulate the dynamic nature of *SPA* by utilizing *AJAX (Asynchronous JavaScript And XML)* to update parts of the page without reloading the whole page.


#figure(
  image("../assets/htmx-vs-spa.png", width: 150mm),
  caption: "HTMX vs SPA Architecture",
)

#figure(
  image("../assets/htmx-flow.png", width: 150mm),
  caption: "HTMX Flow Diagram",
)
== Database Design
To implement such a _performance and data-driven_ application, a database with a _well-designed and optimized schema_ must be designed and implemented. Such a design can be obtained by using the following principles:

+ Database normalization.

+ Database indexing.

+ Soft deletion and archiving.

+ Automatic data integrity with triggers and constraints.

+ Compiled queries and type-safety.

+ Database Transaction Management.

=== Database Normalization

Database normalization is the process of organizing the data in the database and reducing redundancy. Database normalization has 3 main normal forms:


1. *First Normal Form (1NF)*: The first normal form of a database implies that the data in the database be atomic where each field contains only one value.

2. *Second Normal Form (2NF)*: The second normal form of a database is when the database satisfies the condition of 1NF and all the non-key fields depend on the primary key completely.

3. *Third Normal Form (3NF)*: The third normal form of a database is when the database satisfies the conditions of 2NF and there are no transitive dependency among non-key fields.

#figure(
  image("../assets/normalization.png", width: 120mm),
  caption: "NF1, NF2, NF3 Forms",
)

=== Database Indexing
Database indexing is an optimization data structure that improves the database performance. By creating an index for the most frequently used queries, a faster row find is possible @jasmine2010normalizing.

#figure(
  image("../assets/index.webp", width: 180mm),
  caption: "Indexing a database table with a Non-Clustered Index",
)
=== Database Transaction Management
A database has a state of consistency. When a transaction is executed, the database state will temporarily change to an inconsistent state. A transaction manager is responsible for making sure that the database returns to a consistent state after the inconsistent period. By _rolling back_ if the transaction fails, otherwise the transaction will be committed successfully.

#figure(
  image("../assets/tx_state.png", width: 180mm),
  caption: "Transaction State Diagram",
)

#figure(
  image("../assets/tx_manager.webp", width: 150mm),
  caption: "Transaction Manager Diagram",
)

== Compiled queries

In order to ensure the _type-safety_ of a database query, an inline written query is not the best for _type-safety_ and a maintainable code base. A better-suited solution is either to use an *ORM (Object Relational Mapping)* or a *compiled query builder*. The *ORM* approach is _more common_ in modern web development. It gives the user a set of interfaces to implement and generate the desired *SQL query*. *ORMs* are a convenient, _easy-to-use_ approach, but they _may lack control_, _flexibility_, and sometimes a lack of _optimization_ on a certain query. On the other hand, a compiled query builder like *SQLC* is a more _flexible_ and _optimized_ approach; it gives the user more control over the generated *SQL query*, also supports the _repository pattern_, but it may require more effort to implement @sqlc_2026.

#figure(
  [
    ```sql
    -- name: GetUserByID :one
    SELECT * FROM user WHERE id = $1;
    ```
    The compiled code generated by SQLC would look like this:
    ```go
    const getUserByID = `-- name: GetUserByID :one
    SELECT id, username, full_name, password, role FROM users
    WHERE id = $1 AND deleted_at IS NULL
    `

    func (q *Queries) GetUserByID(ctx context.Context, id uuid.UUID) (User, error) {
    	row := q.db.QueryRow(ctx, getUserByID, id)
    	var i User
    	err := row.Scan(
    		&i.ID,
    		&i.Username,
    		&i.FullName,
    		&i.Password,
    		&i.Role,
    	)
    	return i, err
    }
    ```
  ],
  caption: "Compiled queries with SQLC example",
)

== Architectural Patterns in Web Back-End Systems

=== MVC (Model-View-Controller)

The MVC architecture pattern is one of the most common architectural choices in web development. It separates the application into three distinct layers: the model layer, the Controller layer, and the view layer. Each layer is responsible for its own tasks.

+ *The Controller layer* is responsible for handling the user requests and responses; it acts as a mediator between the model and the view layers.

+ *The Model layer*: is responsible for handling the business logic and data manipulation.

+ *The View layer*: is responsible for rendering and presenting the data and the interface to the end user.


#figure(
  image("../assets/mvc.png", width: 100mm),
  caption: "MVC Architecture Diagram",
)
=== Hexagonal Architecture (Ports and Adapters)
*Hexagonal architecture (Ports and Adapters)* is a _software design pattern_ first introduced by *_Alistair Cockburn_* in 2005 @cockburn2005hexagonal. This architectural pattern is based on the idea of separating the application into different decoupled layers: _the core layer_, _the ports layer_, and _the adapters layer_. All layers are wrapped in a _hexagonal shape_ with _the core layer_ in the center that has no dependencies on the other layers, and communicate with the outside world through _well-defined ports and adapters_.

#figure(
  image("../assets/hex_archi.png", width: 100mm),
  caption: "Hexagonal Architecture Diagram",
)

The hexagonal architecture pattern is based on a principle called the *_dependency inversion principle_* @Noback2018, which is the last principle of the *SOLID* principles. That state :

#quote(
  block: true,
)[_High-level modules should not import anything from low-level modules. Both should depend on abstractions (e.g., interfaces)._]


#quote(
  block: true,
)[_Abstractions should not depend on details. Details (concrete implementations) should depend on abstractions._]

== Academic Integrity in Online Examinations

Online examinations are not like traditional _pen-and-paper_ exams, which can be monitored somewhat easily. In a digital online examination _environment_, _academic integrity_ measures must be implemented to ensure the _integrity_ of the examination framework. Some of these measures are :

- Time limits and auto-submission.

- Question randomization.

- Question checksums.

- Audit logs and monitoring.

== Role-Based Access Control (RBAC)

*Role-Based Access Control (RBAC)* is a security model that restricts access to resources based on users' roles. In the examination management system context, *RBAC* can be used to _restrict and separate_ access to different functionalities based on the user role.

#figure(
  image("../assets/rbac.png", width: 100mm),
  caption: "Role-Based Access Control Diagram",
)

In this specific system, the *RBAC* is implemented into distinct layers. The _highest level layer_ is implemented as an _HTTP middleware_ and a *JWT-claims* using the *Go* programming language, which is responsible for checking the user role and allowing or denying access to the requested resource. The other layer is implemented as a _database constraints_ using an _enum type_ to restrict and validate roles `user_role_type`.


```sql
CREATE TYPE user_role_type AS
ENUM ('STUDENT', 'INSTRUCTOR', 'ADMIN');
```

== Time-Bound Exams and Auto-Submission
To preserve academic integrity, the system must implement an auto-submission mechanism that renders the examination session immutable once the duration expires. Ensuring high-fidelity synchronization for this feature necessitates a selection between two distinct architectural communication patterns: Full-Duplex WebSockets or Periodic HTTP Polling.


=== HTTP Polling
HTTP Polling is a traditional communication technique where the client (the browser) periodically sends requests to the server to check for new data or status updates. Unlike WebSockets, which keep a connection open, polling follows the standard "request-response" cycle of the HTTP protocol.


=== WebSockets

The *WebSocket protocol (RFC 6455)* facilitates a persistent, full-duplex communication channel over a single TCP connection. Unlike the standard HTTP request-response paradigm, *WebSockets* bypass the overhead of repeated handshake negotiations, allowing the server to push "*Force-Submit*" events to the client in real-time. This provides a low-latency solution suitable for high-concurrency examination environments where sub-second synchronization is paramount @ws2021.


#figure(
  image("../assets/ws.svg", width: 130mm),
  caption: "WebSockets Vs HTTP Polling",
)

In this system, the *WebSocket* approach is chosen for its superior performance and real-time capabilities, ensuring that all clients receive immediate notifications when the examination duration expires, thus maintaining the integrity of the examination process.

== User Interface
The user interface is designed using *Tailwind CSS* and *TemplUI* component library with *HTMX* for the interactivity. The UI is designed to be _user-friendly_.

The Use of *TemplUI* component library gives the UI a modern and unified look and feel. The library provides a high level and customizable components.

#figure(
  image("../assets/templui.png", width: 130%),
  caption: "TemplUI Component Library",
)
#pagebreak()
== Related Work

=== Moodle
*Moodle* is a free and open-source learning management system (LMS) used for online education and training. It was developed by Martin Dougiamas and first released in 2002. It supports a various functionalities one of them is the quiz management system. The system is a web-based application based on a _PHP_ programming language and a _MySQL_ database. But the system is very resource-intensive with a low performance and a outdated UI.

=== Google Forms and Microsoft Forms
*_Google Forms_* and *_Microsoft Forms_* are two of the most popular e-examination platforms. They both provide a simple and basic _MCQ format_. These platforms are lacking robust features like *_RBAC_*, *_question checksums_*, and *_time-bound exams_*. These platforms are storing all the response data in a cloud-based infrastructure over which the institution has no control.

=== Existing System (College of Computer Science)
The currently used system in the college of _computer science and computer engineering_ is a *CSharp* application that connects with a _PHP_ backend using a _RESTful API_. The system uses an old stack of technologies, with _hardcoded_ admin credentials, no _question checksums_, no _question randomization_, no _images_, no _mathematical notations_, and no _audit logs_. The UI is very outdated and not user-friendly. The system does not utilize any clean architectural patterns, and the code is not maintainable. The system also suffers from _segmentation faults_ and memory leaks. The idea of this project is to solve the existing problems with a new _modern_, _secure_, _scalable_, and _performant_ examination management system.

