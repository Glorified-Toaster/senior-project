#let upm-report(
  // Cover page info
  title: "Title of the Work",
  author: "Author's name",
  supervisor: "Supervisor's Name",
  date: datetime.today(),
  // Acknowledgements and Abstract
  acknowledgements: none,
  abstract-en: none,
  keywords-en: none,
  abstract-ar: none,
  keywords-ar: none,
  title-ar: "",
  author-ar: "",
  supervisor-ar: "",
  university-ar: "",
  college-ar: "",
  department-ar: "",
  date-ar: "",
  // School configuration
  university: "",
  school-name: "",
  school-address: "",
  school-abbr: "",
  report-type: "",
  degree-name: "",
  school-color: rgb(32, 130, 192),
  school-logo: "",
  school-watermark: "",
  // Bibliography configuration
  bibliography-file: none,
  bibliography-style: "ieee",
  supervisorCert: none,
  ExaminationCert: none,
  body,
) = {
  // Cover page
  page(
    margin: 1.5cm,
    paper: "a4",
    numbering: none,
  )[
    #set text(size: 12pt, font: ("Libertinus Serif", "serif"))

    // Header
    #grid(
      columns: (auto, 1fr, auto),
      align: horizon,
      move(image("assets/ce_logo.png", width: 25mm), dy: -16pt),
      align(center)[
        #text(size: 13pt, weight: "bold")[
          Ministry of Higher Education and Scientific Research

          University of Technology

          College of Computer Engineering

          Network Engineering & \ Cyber Security Department
        ]
      ],
      move(image("assets/logoPNG.png", width: 25mm), dy: -16pt),
    )

    #v(1cm)

    // Title
    #align(center)[
      #text(size: 28pt, weight: "bold", fill: school-color)[
        #title
      ]

      #v(2.0cm)

      // Subtitle
      #block(width: 80%)[
        #text(size: 14pt, weight: "bold")[
          Graduation project submitted to the College of Computer Engineering, in partial fulfillment of the requirements for the (B.Sc.) degree in Computer Engineering.
        ]
      ]

      #v(1cm)

      // "By"
      #text(size: 16pt, weight: "bold", fill: school-color)[By]

      #v(0.5cm)

      // Authors
      #let author-list = author.split(",")
      #grid(
        columns: (1fr, 1fr),
        gutter: 1cm,
        ..author-list.map(name => align(center)[#text(size: 14pt, weight: "bold", style: "italic")[#name.trim()]])
      )

      #v(1.5cm)

      // Supervised by
      #text(size: 16pt, weight: "bold", fill: school-color)[Supervised by]

      #v(0.4cm)

      #text(size: 15pt, weight: "bold", style: "italic")[#supervisor]

      #v(1cm)

      // Date
      #text(size: 16pt, weight: "bold", fill: school-color)[
        2025-2026
      ]
    ]
  ]


  // Supervisor page (optional)
  page(
    header: none,
    numbering: none,
  )[
    #align(center + horizon, image("assets/a.jpeg", width: 100%))
  ]

  // Supervisor page (optional)
  if supervisorCert != none {
    page(
      header: none,
      numbering: none,
    )[
      #set text(size: 12pt)

      #block(
        above: 0pt,
        below: 25pt,
        {
          set text(size: 25pt, weight: "regular")
          align(center)[Supervisor's Certification]
          v(-0.9em)
          v(0.5em)
        },
      )

      #set par(justify: true, leading: 0.65em)
      #v(1em)
      #supervisorCert
      #v(32em)
      #text(size: 15pt)[

        Supervisor's signature:

        Supervisor's name:

        Supervisor's Scientific Degree:

        Date:

      ]
    ]
  }

  // Examination page (optional)
  if supervisorCert != none {
    page(
      header: none,
      numbering: none,
    )[
      #set text(size: 12pt)

      #block(
        above: 0pt,
        below: 25pt,
        {
          set text(size: 25pt, weight: "regular")
          align(center)[Examination Committee Certification]
          v(-0.9em)
          v(0.5em)
        },
      )

      #let GridData = text(size: 15pt)[

        Signature:

        Name:

        Scientific Degree:

        Date:
      ]

      #set par(justify: true, leading: 0.65em)
      #v(1em)
      #ExaminationCert
      #v(5em)
      #grid(
        columns: (1.5fr, 1fr),
      )[#GridData \ (Member)][#GridData \ (Member)]
      #v(4em)
      #grid(
        columns: (1.5fr, 1fr),
      )[#GridData \ (Member)][#GridData \ (Chairman)]

      #v(2em)
      #text(size: 15pt)[

        Signature:

        Name:

        Scientific Degree:

        Date:

        Head of Computer Engineering Department
      ]]
  }

  page(
    header: none,
    numbering: none,
  )[
    #set text(size: 12pt)

    #align(center + horizon, block(
      above: 0pt,
      below: 25pt,
      {
        set text(size: 25pt, weight: "regular")
        align(center)[Dedication]
        v(0.5em)
      },
    ))

    #set par(justify: true, leading: 0.65em)
    #text(size: 15pt)[
      This work is dedicated to the academic community in computer engineering and technological innovation to researchers, educators, and students whose continuous pursuit of knowledge drives the evolution of digital systems and intelligent technologies.

      It is further dedicated to those who believe in the transformative power of computing,
      and who contribute to advancing fields such as data science, software engineering, and emerging technologies that shape modern society.

      May this work serve as a modest contribution to the expanding body of scientific knowledge,
      and inspire further exploration, innovation, and excellence in the field of computer systems.
    ]
  ]

  // Acknowledgements page (optional)
  if acknowledgements != none {
    page(
      header: none,
      numbering: none,
    )[
      #set text(size: 12pt)

      #align(center + horizon, block(
        above: 0pt,
        below: 25pt,
        {
          set text(size: 25pt, weight: "regular")
          align(center)[Acknowledgements]
          v(0.5em)
        },
      ))

      #set par(justify: true, leading: 0.65em)
      #acknowledgements
    ]
  }

  // Abstract page (optional)
  if abstract-en != none {
    page(
      header: none,
      numbering: none,
    )[
      #set text(size: 12pt)

      #block(
        above: 0pt,
        below: 25pt,
        {
          set text(size: 30pt, weight: "regular")
          align(left)[Abstract]
          v(-0.9em)
          line(length: 100%, stroke: 0.5pt)
          v(0.5em)
        },
      )

      #set par(justify: true, leading: 0.65em)

      #abstract-en

      #if keywords-en != none [
        #v(1em)
      ]
    ]
  }


  // Table of Contents
  page(
    header: none,
    numbering: none,
  )[
    // Style the outline entries
    #show outline.entry.where(level: 1): it => {
      v(18pt, weak: true)
      text(fill: school-color)[#strong(it)]
    }

    // Style level 2 entries in blue
    #show outline.entry.where(level: 2): it => {
      v(12pt, weak: true)
      text(fill: school-color)[#it]
    }
    // Style level 2 entries in blue
    #show outline.entry.where(level: 3): it => {
      v(11pt, weak: true)
      text()[#it]
    }
    #set text(size: 11pt)

    // Custom title styled like chapter headings
    #block(
      above: 0pt,
      below: 25pt,
      {
        set text(size: 30pt, weight: "regular")
        align(left)[Table of Contents]
        v(-0.9em)
        line(length: 100%, stroke: 0.5pt)
        v(0.5em)
      },
    )

    #set par(leading: 1.8em)

    #outline(
      title: none,
      depth: 3,
      indent: auto,
    )
  ]

  // List of Figures (optional - only shows if document has figures)
  context {
    let figures = query(figure.where(kind: image))
    if figures.len() > 0 {
      page(
        header: none,
        numbering: none,
      )[
        // Style outline entries in blue with spacing
        #show outline.entry: it => {
          v(12pt, weak: true)
          text(fill: school-color)[#it]
        }

        #set text(size: 12pt)

        #block(
          above: 0pt,
          below: 25pt,
          {
            set text(size: 30pt, weight: "regular")
            align(left)[List of Figures]
            v(-0.9em)
            line(length: 100%, stroke: 0.5pt)
            v(0.5em)
          },
        )

        #set par(leading: 1.8em)

        #outline(
          title: none,
          target: figure.where(kind: image),
        )
      ]
    }
  }

  // Page setup for content
  set page(
    paper: "a4",
    margin: (top: 3.5cm, bottom: 3cm, left: 2.5cm, right: 2.5cm),
    header: context {
      let page-num = counter(page).get().first()

      // Don't show header on first page
      if page-num <= 1 {
        return
      }

      // Check if current page has a chapter heading
      let current-page = here().page()
      let headings-on-page = query(heading.where(level: 1)).filter(h => h.location().page() == current-page)

      // Don't show header on pages with chapter headings
      if headings-on-page.len() > 0 {
        return
      }

      // Show header on all other pages
      set text(fill: school-color, size: 10pt)

      let chapter-title = {
        let elems = query(heading.where(level: 1))
        if elems.len() > 0 {
          let relevant = elems.filter(h => h.location().page() <= current-page)
          if relevant.len() > 0 {
            upper(relevant.last().body)
          }
        }
      }

      grid(
        columns: (1fr, auto),
        align: (left, right),
        chapter-title, counter(page).display(),
      )
      line(length: 100%, stroke: 0.5pt + school-color)
    },
  )

  // Reset page counter for main content
  counter(page).update(1)

  // Text setup
  set text(size: 14pt, font: "Libertinus Serif")

  set par(
    leading: 0.65em,
    spacing: 1em,
    justify: true,
    first-line-indent: 0pt,
  )

  // Heading setup
  show heading: set text(font: "Libertinus Serif")

  // Chapter headings (level 1)
  show heading.where(level: 1): it => {
    if it.numbering != none {
      // separator page
      page(header: none, footer: none)[
        #set align(center + horizon)
        #set text(font: "Libertinus Serif", fill: school-color)

        #image("assets/chapter_top.png", width: 40%)
        #v(0.1em)

        #line(length: 80%, stroke: 1pt + school-color)

        #v(0.5em)
        #text(size: 32pt, weight: "bold")[Chapter #counter(heading).display("I")]
        #v(0.1em)
        #text(size: 32pt, weight: "bold")[#it.body]
        #v(0.5em)

        #line(length: 80%, stroke: 1pt + school-color)
        #v(1em)
        // Bottom ornament placeholder
        #image("assets/chapter_down.png", width: 20%)
      ]
    } else {
      // Standard h1 for non-numbered sectons (Bibliography, etc.)
      pagebreak(weak: true)
      set text(size: 30pt, weight: "regular", font: "Libertinus Serif")

      block(
        above: 0pt,
        below: 25pt,
        width: 100%,
        {
          set par(spacing: 0pt)
          align(left)[
            #it.body
            #v(-0.9em)
            #line(length: 100%, stroke: 0.5pt)
          ]
        },
      )
    }
  }

  // Section headings (level 2)
  show heading.where(level: 2): it => {
    set text(size: 20pt, weight: "bold")
    block(
      above: 1.5em,
      below: 1em,
      {
        if it.numbering != none {
          counter(heading).display()
          [. ]
        }
        h(0.6em)
        it.body
      },
    )
  }

  // Subsection headings (level 3)
  show heading.where(level: 3): it => {
    set text(size: 16pt, weight: "bold")
    block(
      above: 1.2em,
      below: 0.8em,
      {
        if it.numbering != none {
          counter(heading).display()
          [. ]
        }
        h(0.6em)
        it.body
      },
    )
  }

  set heading(numbering: "1.1")

  body

  // Bibliography section
  if bibliography-file != none {
    pagebreak()

    // Style the bibliography heading
    show bibliography: set heading(numbering: none)

    bibliography(
      bibliography-file,
      title: [Bibliography],
      style: bibliography-style,
    )
  }

  // Abstract page (optional)
  if abstract-en != none {
    page(
      header: none,
      numbering: none,
    )[
      #set text(size: 12pt)

      #block(
        above: 0pt,
        below: 25pt,
        {
          set text(size: 30pt, weight: "regular")
          align(right)[الخلاصة]
          v(-0.6em)
          line(length: 100%, stroke: 0.5pt)
          v(0.5em)
        },
      )

      #set par(justify: true, leading: 0.65em)

      #abstract-ar

      #if keywords-en != none [
        #v(1em)
      ]
    ]
  }

  // Arabic Title Page (optional)
  if title-ar != "" {
    page(
      margin: 1.5cm,
      paper: "a4",
      header: none,
      numbering: none,
    )[
      #set text(size: 12pt, font: ("Noto Naskh Arabic", "Libertinus Serif"))
      #set align(center)

      // Header
      #grid(
        columns: (auto, 1fr, auto),
        align: horizon,
        move(image("assets/logoPNG.png", width: 25mm), dy: -10pt),
        align(center)[
          #text(size: 14pt, weight: "bold")[
            وزارة التعليم العالي والبحث العلمي

            #university-ar

            #college-ar

            #department-ar
          ]
        ],
        move(image("assets/ce_logo.png", width: 25mm), dy: -10pt),
      )

      #v(2cm)

      #text(size: 26pt, weight: "bold", fill: school-color)[
        #title-ar
      ]

      #v(1cm)

      #block(width: 85%)[
        #text(size: 15pt, weight: "bold")[
          مشروع مقدم إلى #college-ar في #university-ar كجزء من متطلبات نيل درجة البكالوريوس في هندسة الحاسوب
        ]
      ]

      #v(1.5cm)

      #text(size: 18pt, weight: "bold", fill: school-color)[بإشراف]
      #v(0.5cm)
      #text(size: 16pt, weight: "bold")[#supervisor-ar]

      #v(1.5cm)

      #text(size: 18pt, weight: "bold", fill: school-color)[إعداد الطلبة]
      #v(0.8cm)

      #let author-list-ar = author-ar.split("،")
      #grid(
        columns: 1fr * calc.min(author-list-ar.len(), 2),
        gutter: 1cm,
        ..author-list-ar.map(name => align(center)[#text(size: 16pt, weight: "bold")[#name.trim()]])
      )

      #v(1fr)

      #text(size: 16pt, weight: "bold", fill: school-color)[
        #date-ar
      ]
    ]
  }
}
