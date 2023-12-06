#let datetimeNow = datetime.today();

#let fliPar(body)={
    par(
        [#h(2em)#body]
    )
}

#let ZRYGenericDocs(
    authors: ("ZRY",),
    title: "标题",
    subTitle: "副标题",
    titleDesc: "标题说明",
    date: datetimeNow.display("[year]-[month]-[day]"),
    docVer: "V0.1",
    pageSettings: (
        paper: "a4",
        leftMargin: 1.25in,
        topMargin: 1in,
        rightMargin: 1.25in,
        bottomMargin: 1in
    ),
    fontSets: (
        sans: "Noto Sans CJK SC",
        //sans: "Source Han Sans",
        serif: "Noto Serif CJK SC",
        //serif: "Source Han Serif VF"
    ),
    coverStyle:(
        titleSize: 36pt,
        titleWeight: 600,
        subTitleSize: 24pt,
        subTitleWeight: 600,
        titleDescSize: 16pt,
        titleDescWeight: 400,
        dateSize: 16pt,
        dateWeight: 400,
        authorsSize: 16pt,
        authorsWeight: 600,
        docVerSize: 16pt,
        docVerWeight: 400,
        topMargin: 72pt,
        subTitleMargin: 72pt,
        dateMargin: 16pt,
        authorsMargin: 16pt,
        titleDescMargin: 72pt,
        authorsJoinSymbol: "，"
    ),
    textStyle:(
        h1:(
            sans: false,
            size: 22pt,
            weight: 800,
            style: "italic",
            stretch: 100%,
            fill: rgb("#000000"),
            tracking: 0pt,
            spacing: 100%,
            overhang: true,
            vspaceBottom: 12pt,
        ),
        h2:(
            sans: true,
            size: 16pt,
            weight: 400,
            style: "normal",
            stretch: 100%,
            fill: rgb("#000000"),
            tracking: 0pt,
            spacing: 100%,
            overhang: true,
            vspaceBottom: 6pt,
        ),
        h3: (
            sans: false,
            size: 14pt,
            weight: 900,
            style: "normal",
            stretch: 100%,
            fill: rgb("#000000"),
            tracking: 0pt,
            spacing: 100%,
            overhang: true,
            vspaceBottom: 12pt,
        ),
        h4: (
            sans: true,
            size: 14pt,
            weight: 500,
            style: "normal",
            stretch: 100%,
            fill: rgb("#000000"),
            tracking: 0pt,
            spacing: 100%,
            overhang: true,
            vspaceBottom: 12pt,
        ),
        h5: (
            sans: true,
            size: 14pt,
            weight: 300,
            style: "italic",
            stretch: 100%,
            fill: rgb("#000000"),
            tracking: 0pt,
            spacing: 100%,
            overhang: true,
            vspaceBottom: 8pt,
        ),
        text: (
            size: 12pt,
            weight: 200,
        ),
    ),
    body
) = {
    // Page Setting
    set page(
        paper: pageSettings.paper,
        margin: (
            left: pageSettings.leftMargin,
            right: pageSettings.rightMargin,
            top: pageSettings.topMargin,
            bottom: pageSettings.bottomMargin
        )
    )
    // Document Meta Setting
    set document(
        title: title,
        author: authors
    )

    // ==== Cover Begin ====
    v(coverStyle.topMargin)
    // Title
    align(center, text(
        font: fontSets.sans,
        weight: coverStyle.titleWeight,
        size: coverStyle.titleSize,
        title
    ))
    v(coverStyle.subTitleMargin)
    // subTitle
    align(center, text(
        font: fontSets.sans,
        weight: coverStyle.subTitleWeight,
        size: coverStyle.subTitleSize,
        subTitle
    ))

    align(bottom, [
        // titleDesc
        #align(center, text(
            font: fontSets.sans,
            weight: coverStyle.titleDescWeight,
            size: coverStyle.titleDescSize,
            titleDesc
        ))
        #v(coverStyle.titleDescMargin)
        // authors
        #align(center, text(
            font: fontSets.sans,
            weight: coverStyle.authorsWeight,
            size: coverStyle.authorsSize,
            authors.join(coverStyle.authorsJoinSymbol)
        ))
        #v(coverStyle.authorsMargin)
        // date
        #align(center, text(
            font: fontSets.sans,
            weight: coverStyle.dateWeight,
            size: coverStyle.dateSize,
            date
        ))
        #v(coverStyle.dateMargin)
        // docVer
        #align(center, text(
            font: fontSets.sans,
            weight: coverStyle.docVerWeight,
            size: coverStyle.docVerSize,
            docVer
        ))
    ])
    // ==== Cover End ====
    pagebreak()

    // Default Font Setting
    set text(
        font: fontSets.serif,
        weight: textStyle.text.weight,
        size: textStyle.text.size,
        lang: "zh",
    )
    // Paragraph Setting
    /*
    show par: set block(spacing: 0.65em)
    set par(
        leading: 0.65em,
        //first-line-indent: 2em,
        justify: true,
    )
    */

    // ==== Heading Style Setting Begin ====
    set heading(numbering: "1.1.1.1.1")
    // H1
    show heading.where(level: 1): it => block(width: 100%)[
        #let fontFamily = if textStyle.h1.sans {
            fontSets.sans
        } else {
            fontSets.serif
        }
        #text(
            font: fontFamily,
            size: textStyle.h1.size,
            weight: textStyle.h1.weight,
            style: textStyle.h1.style,
            stretch: textStyle.h1.stretch,
            fill: textStyle.h1.fill,
            tracking: textStyle.h1.tracking,
            spacing: textStyle.h1.spacing,
            overhang: textStyle.h1.overhang,
            counter(heading).display() + h(1em) + it.body
        )
        #v(textStyle.h1.vspaceBottom)
    ]
    // H2
    show heading.where(level: 2): it => block(width: 100%)[
        #let fontFamily = if textStyle.h2.sans {
            fontSets.sans
        } else {
            fontSets.serif
        }
        #text(
            font: fontFamily,
            size: textStyle.h2.size,
            weight: textStyle.h2.weight,
            style: textStyle.h2.style,
            stretch: textStyle.h2.stretch,
            fill: textStyle.h2.fill,
            tracking: textStyle.h2.tracking,
            spacing: textStyle.h2.spacing,
            overhang: textStyle.h2.overhang,
            counter(heading).display() + h(1em) + it.body
        )
        #v(textStyle.h2.vspaceBottom)
    ]
    // H3
    show heading.where(level: 3): it => block(width: 100%)[
        #let fontFamily = if textStyle.h3.sans {
            fontSets.sans
        } else {
            fontSets.serif
        }
        #text(
            font: fontFamily,
            size: textStyle.h3.size,
            weight: textStyle.h3.weight,
            style: textStyle.h3.style,
            stretch: textStyle.h3.stretch,
            fill: textStyle.h3.fill,
            tracking: textStyle.h3.tracking,
            spacing: textStyle.h3.spacing,
            overhang: textStyle.h3.overhang,
            counter(heading).display() + h(1em) + it.body
        )
        #v(textStyle.h3.vspaceBottom)
    ]
    // H4
    show heading.where(level: 4): it => block(width: 100%)[
        #let fontFamily = if textStyle.h4.sans {
            fontSets.sans
        } else {
            fontSets.serif
        }
        #text(
            font: fontFamily,
            size: textStyle.h4.size,
            weight: textStyle.h4.weight,
            style: textStyle.h4.style,
            stretch: textStyle.h4.stretch,
            fill: textStyle.h4.fill,
            tracking: textStyle.h4.tracking,
            spacing: textStyle.h4.spacing,
            overhang: textStyle.h4.overhang,
            counter(heading).display() + h(1em) + it.body
        )
        #v(textStyle.h4.vspaceBottom)
    ]
    // H5
    show heading.where(level: 5): it => block(width: 100%)[
        #let fontFamily = if textStyle.h5.sans {
            fontSets.sans
        } else {
            fontSets.serif
        }
        #text(
            font: fontFamily,
            size: textStyle.h5.size,
            weight: textStyle.h5.weight,
            style: textStyle.h5.style,
            stretch: textStyle.h5.stretch,
            fill: textStyle.h5.fill,
            tracking: textStyle.h5.tracking,
            spacing: textStyle.h5.spacing,
            overhang: textStyle.h5.overhang,
            counter(heading).display() + h(1em) + it.body
        )
        #v(textStyle.h5.vspaceBottom)
    ]
    // ==== Title Style Setting End ====

    // ==== Body Begin ====
    body
    // ==== Body End ====
}