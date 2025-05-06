import { CompletionContext } from "@codemirror/autocomplete";

export function myCompletions(context: CompletionContext) {
    let word = context.matchBefore(/\w*/);
    if (!word || (word.from == word.to && !context.explicit)) return null;
  
    return {
      from: word.from,
      options: [
        { label: "begin", type: "keyword", 
            apply: "begin{document}", info: "begin document" },
            { label: "end", type: "keyword", 
            apply: "end{document}", info: "end document" },
            { label: "section", type: "keyword", 
            apply: "section{}", info: "section" },
            { label: "subsection", type: "keyword", 
            apply: "subsection{}", info: "subsection" },
            { label: "subsubsection", type: "keyword", 
            apply: "subsubsection{}", info: "subsubsection" },
            { label: "itemize", type: "keyword", 
            apply: "begin{itemize}\\item\\end{itemize}", info: "itemize" },
            { label: "enumerate", type: "keyword", 
            apply: "begin{enumerate}\\item\\end{enumerate}", info: "enumerate"},
            { label: "description", type: "keyword", 
            apply: "begin{description}\\item[]\\end{description}", info: "description" },
        { label: "figure", type: "keyword", 
            apply: "begin{figure}\\includegraphics{}\\caption{}\\end{figure}", info: "figure" },
        { label: "table", type: "keyword", 
            apply: "begin{table}\\begin{tabular}{}\\end{tabular}\\caption{}\\end{table}", info: "table" },
        { label: "align", type: "keyword", 
            apply: "begin{align}\\end{align}", info: "align" },
        { label: "equation", type: "keyword", 
            apply: "begin{equation}\\end{equation}", info: "equation" },
        { label: "item", type: "keyword", 
            apply: "item{}", info: "item" },
        { label: "textbf{}", type: "keyword", 
            apply: "textbf{}", info: "Bold text" },
        { label: "textit{}", type: "keyword", 
            apply: "textit{}", info: "\\textit{}" },
        { label: "underline{}", type: "keyword", 
            apply: "underline{}", info: "\\underline{}" },
        { label: "newline{}", type: "keyword", 
            apply: "newline{}", info: "\\newline{}" }
      ]
    };
  }