// ZRY Generic Docs Template For Typst.
let ZRYGenericDocs(
    author:"",
    title: "",
    subTitle: "",
    titleDesc: "",
    date: "",
    docVer: "",

    body
) = {
    set page(
        paper: "a4",
        margin: (
            left: 1.25in,
            top: 1in,
            right: 1.25in,
            bottom: 1in
        )
    )
    TEST
    body
}