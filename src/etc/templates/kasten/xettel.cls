\NeedsTeXFormat{LaTeX2e}
\ProvidesClass{xettel}[2024/09/09 XK Zettel Class]

\LoadClass[class=article, crop=false]{standalone}

\RequirePackage{../preamble}
\RequirePackage[backend=biber]{biblatex}
\RequirePackage{hyperref}

\addbibresource{../zettelkasten.bib}

\DeclareBibliographyDriver{zettel}{%
	\printfield{title}%
	\setunit{\addspace}%
	,\printnames{author}%
}

\AtBeginDocument{}
\AtEndDocument{
	\printbibliography
}

\endinput
