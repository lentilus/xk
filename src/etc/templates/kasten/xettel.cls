\NeedsTeXFormat{LaTeX2e}
\ProvidesClass{xettel}[2024/09/09 XK Zettel Class]

% \RequirePackage{mdframed}

\LoadClass[class=article, crop=false]{standalone}

\RequirePackage{../testingpreamble}
% \RequirePackage{currfile}

\AtBeginDocument {
	\begin{mdframed}
		}

		\AtEndDocument {
	\end{mdframed}
}

\endinput
