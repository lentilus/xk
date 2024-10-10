\NeedsTeXFormat{LaTeX2e}
\ProvidesClass{xettel}[2024/09/09 XK Zettel Class]

\LoadClass[class=article, crop=false]{standalone}

\RequirePackage{../preamble}

\AtBeginDocument{\begin{zettel} }
		\AtEndDocument{\end{zettel}}

\endinput
