= Results and Evaluation

== System Architecture Evaluation

=== Comparison: Monolithic MVC vs. Hexagonal Architecture
This section evaluates this system's architecture (*Hexagonal Architecture*) against the monolithic *MVC* architecture like the one used in *Moodle*.

#figure(
  table(
    columns: 3,
    [Feature], [MVC], [Hexagonal Architecture],
    [Coupling],
    [
 The model is tightly coupled to the view and controller.
    ],
    [
 The core domain is decoupled from the dependencies of the other system layers.
    ],

    [Testability],
    [
 Hard to test due to the tight coupling between layers. If a function that requires a database connection is tested, it must have a database connection.
    ],
    [
 Easy to test due to the separation of concerns and layers. A mock database connection can be used to test the same function without a database connection.
    ],

    [Replaceability],
    [
 Changing the database or web framework requires a system-wide rewrite.
    ],
    [
 Changing the database or web framework is easy and requires rewriting only the adapters without touching the core domain of the system.
    ],

    [Complexity], [Lower initial complexity], [Higher initial complexity],
    [Framework \ lock-in],
    [The framework is tightly coupled to the system, making it difficult to switch to a different framework.],
    [The framework is decoupled from the core domain of the system, and only located in the adapters layer],
  ),
)

The hexagonal architecture approach seems at first glance to be more complex and more verbose than the traditional *MVC*. But, it provides a better architectural solution when the codebase grows, and the system becomes more _complex_.

== Performance Evaluation
The system performance was evaluated using:

- *Grafana K6*: A load testing tool for *APIs*, *microservices*, and *websites*.

The hardware used for the evaluation is:

- CPU: AMD Ryzen 7 4800H (16 Threads) `@` 4.30GHz
- RAM: 16GB
- GPU: NVIDIA RTX 3060
- OS: Arch Linux 64-bit (Kernel 7.0.0-1-cachyos)
- Docker Engine: 29.4.0

The *k6* script is as follows:

```JavaScript
import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
 vus: 100,
 duration: '30s',
};

export default function () {
 const res = http.get('https://localhost:8443/health', {
 tlsInsecureSkipVerify: true,
 });
 check(res, {
 'status is 200': (r) => r.status === 200,
 });
 sleep(1);
}
```
=== Response Time
A different request endpoint was tested and aggregated into a *_CSV_* file, then visualized.
#figure(
  image("../assets/performance_report.png"),
  caption: [Response time of the system under load using *k6* and *Matplotlib*],
)

=== Memory Usage
#figure(
  image("../assets/memory_normalized.png"),
  caption: [Memory usage of the system under load using *k6*],
)

== Comparison with Moodle, and MCQUOT

#figure(
  table(
    columns: 4,
    [Metric], [Moodle], [MCQUOT], [This System],
    [Response Time \ (Average)], [700ms], [660ms], [3ms],
    [Memory Usage \ (Idle)], [97MB], [83MB], [35MB],
    [Memory Usage \ (Peak)], [1GB to every 15 users], [1GB to every 20 users], [52MB to every 100 users],
    [Startup time], [1-2s], [1-2s], [`<100 ms`],
    [Language], [PHP], [PHP], [Go],
    [Concurrency model], [Multi-process], [Multi-process], [Goroutines \ (lightweight; M:N \ scheduling)],
    [Type safety], [Dynamic typing], [Dynamic typing], [Strong static typing],
    [Binary size],
    [
 N/A],
    [N/A],
    [60MB \ (self-contained \ binary)],

    [Source code size], [434.8 MB], [422.9 MB], [58.4 MB],
  ),
)

== Solved Problems of the Previous MCQUOT System
The previous system was suffering from a lot of _architectural_, _performance_, _security_, and _functional_ problems. Some of these functional problems are:

=== Database question duplication
The previous system had no mechanism to prevent _duplication_ of questions in a single exam context. So the same question could appear multiple times in the same exam.

To prevent this, we added a constraint that computes the `MD5` hash of the question text and the exam `uuid` to ensure that there is no duplication of questions in the same exam.


=== Database question randomization
The previous system had no mechanism to _randomize_ the questions. To _randomize_ the questions in a single exam attempt, we added a hash of the exam's `uuid` with the exam_attempt's `uuid` to generate a random seed for the random number generator to order the questions on it.

=== Question limitation
The previous system had only one type of question, and it was simply the text question type. This system supports multiple types of questions, such as:

==== Text questions
#figure(
  image("../assets/text_question.png"),
  caption: [Text question type],
)
==== Code highlighted questions
#figure(
  image("../assets/code_question.png"),
  caption: [Code highlighted question type],
)
==== Image questions
#figure(
  image("../assets/image_question.png"),
  caption: [Image question type],
)


=== Manual Data Entry
The previous system had no mechanism to automate the data entry process. This system is solving this problem by implementing a *CSV* _import/export_ feature to add questions, exams, users, and subjects to the system.