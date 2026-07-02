package fzf

import (
	"bytes"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/junegunn/fzf/src/algo"
)

type IconEntry struct {
	Key   string
	Icon  string
	Color int
}

var FilenameIcons = []IconEntry{
	{".babelrc", "", 143},
	{".bash_profile", "", 113},
	{".bashrc", "", 113},
	{".clang-format", "", 102},
	{".clang-tidy", "", 102},
	{".clangd", "", 102},
	{".codespellrc", "󰓆", 77},
	{".condarc", "", 70},
	{".dockerignore", "󰡨", 68},
	{".ds_store", "", 59},
	{".editorconfig", "", 224},
	{".env", "", 185},
	{".eslintignore", "", 55},
	{".eslintrc", "", 55},
	{".git-blame-ignore-revs", "", 166},
	{".gitattributes", "", 166},
	{".gitconfig", "", 166},
	{".gitignore", "", 166},
	{".gitlab-ci.yml", "", 166},
	{".gitmodules", "", 166},
	{".gtkrc-2.0", "", 231},
	{".gvimrc", "", 29},
	{".justfile", "", 102},
	{".luacheckrc", "", 39},
	{".luaurc", "", 39},
	{".mailmap", "󰊢", 166},
	{".nanorc", "", 54},
	{".npmignore", "", 161},
	{".npmrc", "", 161},
	{".nuxtrc", "󱄆", 36},
	{".nvmrc", "", 71},
	{".pnpmfile.cjs", "", 178},
	{".pre-commit-config.yaml", "󰛢", 178},
	{".prettierignore", "", 68},
	{".prettierrc", "", 68},
	{".prettierrc.cjs", "", 68},
	{".prettierrc.js", "", 68},
	{".prettierrc.json", "", 68},
	{".prettierrc.json5", "", 68},
	{".prettierrc.mjs", "", 68},
	{".prettierrc.toml", "", 68},
	{".prettierrc.yaml", "", 68},
	{".prettierrc.yml", "", 68},
	{".pylintrc", "", 102},
	{".settings.json", "", 97},
	{".srcinfo", "󰣇", 32},
	{".vimrc", "", 29},
	{".xauthority", "", 166},
	{".xinitrc", "", 166},
	{".xresources", "", 166},
	{".xsession", "", 166},
	{".zprofile", "", 113},
	{".zshenv", "", 113},
	{".zshrc", "", 113},
	{"_gvimrc", "", 29},
	{"_vimrc", "", 29},
	{"authors", "", 141},
	{"authors.txt", "", 141},
	{"bitbucket-pipelines.yml", "󰂨", 33},
	{"brewfile", "", 88},
	{"bspwmrc", "", 235},
	{"build", "", 113},
	{"build.gradle", "", 24},
	{"build.zig.zon", "", 178},
	{"bun.lock", "", 188},
	{"bun.lockb", "", 188},
	{"cantorrc", "", 38},
	{"checkhealth", "󰓙", 110},
	{"cmakelists.txt", "", 188},
	{"code_of_conduct", "", 161},
	{"code_of_conduct.md", "", 161},
	{"commit_editmsg", "", 166},
	{"commitlint.config.js", "󰜘", 30},
	{"commitlint.config.ts", "󰜘", 30},
	{"compose.yaml", "󰡨", 68},
	{"compose.yml", "󰡨", 68},
	{"config", "", 102},
	{"containerfile", "󰡨", 68},
	{"copying", "", 143},
	{"copying.lesser", "", 143},
	{"directory.build.props", "", 39},
	{"directory.build.targets", "", 39},
	{"directory.packages.props", "", 39},
	{"docker-compose.yaml", "󰡨", 68},
	{"docker-compose.yml", "󰡨", 68},
	{"dockerfile", "󰡨", 68},
	{"dune", "", 94},
	{"eslint.config.cjs", "", 55},
	{"eslint.config.js", "", 55},
	{"eslint.config.mjs", "", 55},
	{"eslint.config.ts", "", 55},
	{"ext_typoscript_setup.txt", "", 208},
	{"favicon.ico", "", 143},
	{"fp-info-cache", "", 231},
	{"fp-lib-table", "", 231},
	{"freecad.conf", "", 131},
	{"gemfile", "", 88},
	{"gnumakefile", "", 102},
	{"go.mod", "", 38},
	{"go.sum", "", 38},
	{"go.work", "", 38},
	{"gradle-wrapper.properties", "", 24},
	{"gradle.properties", "", 24},
	{"gradlew", "", 24},
	{"groovy", "", 66},
	{"gruntfile.babel.js", "", 173},
	{"gruntfile.coffee", "", 173},
	{"gruntfile.js", "", 173},
	{"gruntfile.ts", "", 173},
	{"gtkrc", "", 231},
	{"gulpfile.babel.js", "", 167},
	{"gulpfile.coffee", "", 167},
	{"gulpfile.js", "", 167},
	{"gulpfile.ts", "", 167},
	{"hypridle.conf", "", 37},
	{"hyprland.conf", "", 37},
	{"hyprlandd.conf", "", 37},
	{"hyprlock.conf", "", 37},
	{"hyprpaper.conf", "", 37},
	{"hyprsunset.conf", "", 37},
	{"i18n.config.js", "󰗊", 103},
	{"i18n.config.ts", "󰗊", 103},
	{"i3blocks.conf", "", 188},
	{"i3status.conf", "", 188},
	{"index.theme", "", 36},
	{"ionic.config.json", "", 68},
	{"jenkinsfile", "", 167},
	{"justfile", "", 102},
	{"kalgebrarc", "", 38},
	{"kdeglobals", "", 38},
	{"kdenlive-layoutsrc", "", 110},
	{"kdenliverc", "", 110},
	{"kritadisplayrc", "", 170},
	{"kritarc", "", 170},
	{"license", "", 179},
	{"license.md", "", 179},
	{"lxde-rc.xml", "", 245},
	{"lxqt.conf", "", 32},
	{"makefile", "", 102},
	{"mix.lock", "", 139},
	{"mpv.conf", "", 53},
	{"next.config.cjs", "", 231},
	{"next.config.js", "", 231},
	{"next.config.ts", "", 231},
	{"node_modules", "", 161},
	{"nuxt.config.cjs", "󱄆", 36},
	{"nuxt.config.js", "󱄆", 36},
	{"nuxt.config.mjs", "󱄆", 36},
	{"nuxt.config.ts", "󱄆", 36},
	{"package-lock.json", "", 88},
	{"package.json", "", 161},
	{"pkgbuild", "", 32},
	{"platformio.ini", "", 172},
	{"playwright.config.cjs", "", 35},
	{"playwright.config.cts", "", 35},
	{"playwright.config.js", "", 35},
	{"playwright.config.mjs", "", 35},
	{"playwright.config.mts", "", 35},
	{"playwright.config.ts", "", 35},
	{"pnpm-lock.yaml", "", 178},
	{"pnpm-workspace.yaml", "", 178},
	{"pom.xml", "", 88},
	{"prettier.config.cjs", "", 68},
	{"prettier.config.js", "", 68},
	{"prettier.config.mjs", "", 68},
	{"prettier.config.ts", "", 68},
	{"prisma.config.mts", "", 68},
	{"prisma.config.ts", "", 68},
	{"procfile", "", 139},
	{"prusaslicer.ini", "", 172},
	{"prusaslicergcodeviewer.ini", "", 172},
	{"py.typed", "", 214},
	{"qtproject.conf", "", 77},
	{"rakefile", "", 88},
	{"readme", "󰂺", 254},
	{"readme.md", "󰂺", 254},
	{"rmd", "", 73},
	{"robots.txt", "󰚩", 66},
	{"security", "󰒃", 145},
	{"security.md", "󰒃", 145},
	{"settings.gradle", "", 24},
	{"svelte.config.js", "", 202},
	{"sxhkdrc", "", 235},
	{"sym-lib-table", "", 231},
	{"tailwind.config.js", "󱏿", 38},
	{"tailwind.config.mjs", "󱏿", 38},
	{"tailwind.config.ts", "󱏿", 38},
	{"tmux.conf", "", 34},
	{"tmux.conf.local", "", 34},
	{"tsconfig.json", "", 73},
	{"unlicense", "", 179},
	{"vagrantfile", "", 27},
	{"vercel.json", "", 231},
	{"vite.config.cjs", "", 214},
	{"vite.config.cts", "", 214},
	{"vite.config.js", "", 214},
	{"vite.config.mjs", "", 214},
	{"vite.config.mts", "", 214},
	{"vite.config.ts", "", 214},
	{"vitest.config.cjs", "", 106},
	{"vitest.config.cts", "", 106},
	{"vitest.config.js", "", 106},
	{"vitest.config.mjs", "", 106},
	{"vitest.config.mts", "", 106},
	{"vitest.config.ts", "", 106},
	{"vlcrc", "󰕼", 172},
	{"webpack", "󰜫", 73},
	{"weston.ini", "", 214},
	{"workspace", "", 113},
	{"wrangler.jsonc", "", 172},
	{"wrangler.toml", "", 172},
	{"xdph.conf", "", 37},
	{"xmobarrc", "", 167},
	{"xmobarrc.hs", "", 167},
	{"xmonad.hs", "", 167},
	{"xorg.conf", "", 166},
	{"xsettingsd.conf", "", 166},
}

var ExtensionIcons = []IconEntry{
	{".babelrc", "", 143},
	{".bash_profile", "", 113},
	{".bashrc", "", 113},
	{".clang-format", "", 102},
	{".clang-tidy", "", 102},
	{".clangd", "", 102},
	{".codespellrc", "󰓆", 77},
	{".condarc", "", 70},
	{".dockerignore", "󰡨", 68},
	{".ds_store", "", 59},
	{".editorconfig", "", 224},
	{".env", "", 185},
	{".eslintignore", "", 55},
	{".eslintrc", "", 55},
	{".git-blame-ignore-revs", "", 166},
	{".gitattributes", "", 166},
	{".gitconfig", "", 166},
	{".gitignore", "", 166},
	{".gitlab-ci.yml", "", 166},
	{".gitmodules", "", 166},
	{".gtkrc-2.0", "", 231},
	{".gvimrc", "", 29},
	{".justfile", "", 102},
	{".luacheckrc", "", 39},
	{".luaurc", "", 39},
	{".mailmap", "󰊢", 166},
	{".nanorc", "", 54},
	{".npmignore", "", 161},
	{".npmrc", "", 161},
	{".nuxtrc", "󱄆", 36},
	{".nvmrc", "", 71},
	{".pnpmfile.cjs", "", 178},
	{".pre-commit-config.yaml", "󰛢", 178},
	{".prettierignore", "", 68},
	{".prettierrc", "", 68},
	{".prettierrc.cjs", "", 68},
	{".prettierrc.js", "", 68},
	{".prettierrc.json", "", 68},
	{".prettierrc.json5", "", 68},
	{".prettierrc.mjs", "", 68},
	{".prettierrc.toml", "", 68},
	{".prettierrc.yaml", "", 68},
	{".prettierrc.yml", "", 68},
	{".pylintrc", "", 102},
	{".settings.json", "", 97},
	{".srcinfo", "󰣇", 32},
	{".vimrc", "", 29},
	{".xauthority", "", 166},
	{".xinitrc", "", 166},
	{".xresources", "", 166},
	{".xsession", "", 166},
	{".zprofile", "", 113},
	{".zshenv", "", 113},
	{".zshrc", "", 113},
	{"1", "", 102},
	{"3gp", "", 172},
	{"3mf", "󰆧", 244},
	{"7z", "", 178},
	{"Dockerfile", "󰡨", 68},
	{"R", "󰟔", 31},
	{"_gvimrc", "", 29},
	{"_vimrc", "", 29},
	{"a", "", 188},
	{"aac", "", 39},
	{"ada", "", 75},
	{"adb", "", 75},
	{"ads", "", 139},
	{"ai", "", 143},
	{"aif", "", 39},
	{"aiff", "", 39},
	{"alma", "", 203},
	{"alpine", "", 24},
	{"android", "", 71},
	{"aosc", "", 124},
	{"ape", "", 39},
	{"apk", "", 71},
	{"apl", "", 35},
	{"app", "", 124},
	{"apple", "", 145},
	{"applescript", "", 102},
	{"arch", "󰣇", 32},
	{"archcraft", "", 109},
	{"archlabs", "", 59},
	{"arcolinux", "", 104},
	{"artix", "", 74},
	{"asc", "󰦝", 66},
	{"asm", "", 31},
	{"ass", "󰨖", 214},
	{"astro", "", 168},
	{"authors", "", 141},
	{"authors.txt", "", 141},
	{"avi", "", 172},
	{"avif", "", 139},
	{"awesomewm", "", 60},
	{"awk", "", 59},
	{"azcli", "", 32},
	{"bak", "󰁯", 102},
	{"bash", "", 113},
	{"bat", "", 148},
	{"bazel", "", 113},
	{"bib", "󱉟", 143},
	{"bicep", "", 73},
	{"bicepparam", "", 139},
	{"biglinux", "", 37},
	{"bin", "", 124},
	{"bitbucket-pipelines.yml", "󰂨", 33},
	{"blade.php", "", 167},
	{"blend", "󰂫", 172},
	{"blp", "󰺾", 68},
	{"bmp", "", 139},
	{"bqn", "", 35},
	{"brep", "󰻫", 101},
	{"brewfile", "", 88},
	{"bspwm", "", 239},
	{"bspwmrc", "", 235},
	{"budgie", "", 59},
	{"build", "", 113},
	{"build.gradle", "", 24},
	{"build.zig.zon", "", 178},
	{"bun.lock", "", 188},
	{"bun.lockb", "", 188},
	{"bz", "", 178},
	{"bz2", "", 178},
	{"bz3", "", 178},
	{"bzl", "", 113},
	{"c", "", 75},
	{"c++", "", 168},
	{"cache", "", 231},
	{"cantorrc", "", 38},
	{"cast", "", 172},
	{"cbl", "", 25},
	{"cc", "", 168},
	{"ccm", "", 168},
	{"centos", "", 132},
	{"cfc", "", 37},
	{"cfg", "", 102},
	{"cfm", "", 37},
	{"checkhealth", "󰓙", 110},
	{"cinnamon", "", 172},
	{"cjs", "", 143},
	{"clj", "", 107},
	{"cljc", "", 107},
	{"cljd", "", 73},
	{"cljs", "", 73},
	{"cmake", "", 188},
	{"cmakelists.txt", "", 188},
	{"cob", "", 25},
	{"cobol", "", 25},
	{"code_of_conduct", "", 161},
	{"code_of_conduct.md", "", 161},
	{"coffee", "", 143},
	{"commit_editmsg", "", 166},
	{"commitlint.config.js", "󰜘", 30},
	{"commitlint.config.ts", "󰜘", 30},
	{"compose.yaml", "󰡨", 68},
	{"compose.yml", "󰡨", 68},
	{"conda", "", 70},
	{"conf", "", 102},
	{"config", "", 102},
	{"config.ru", "", 88},
	{"containerfile", "󰡨", 68},
	{"copying", "", 143},
	{"copying.lesser", "", 143},
	{"cow", "󰆚", 94},
	{"cp", "", 73},
	{"cpp", "", 73},
	{"cppm", "", 73},
	{"cpy", "", 25},
	{"cr", "", 251},
	{"crdownload", "", 79},
	{"crystallinux", "", 129},
	{"cs", "󰌛", 64},
	{"csh", "", 59},
	{"cshtml", "󱦗", 56},
	{"cson", "", 143},
	{"csproj", "󰪮", 56},
	{"css", "", 97},
	{"csv", "", 113},
	{"cts", "", 73},
	{"cu", "", 113},
	{"cue", "󰲹", 175},
	{"cuh", "", 139},
	{"cxx", "", 73},
	{"cxxm", "", 73},
	{"d", "", 130},
	{"d.ts", "", 173},
	{"dart", "", 25},
	{"db", "", 188},
	{"dconf", "", 231},
	{"debian", "", 124},
	{"deepin", "", 38},
	{"desktop", "", 60},
	{"devuan", "", 59},
	{"diff", "", 59},
	{"directory.build.props", "", 39},
	{"directory.build.targets", "", 39},
	{"directory.packages.props", "", 39},
	{"dll", "", 52},
	{"doc", "󰈬", 25},
	{"docker-compose.yaml", "󰡨", 68},
	{"docker-compose.yml", "󰡨", 68},
	{"dockerfile", "󰡨", 68},
	{"dockerignore", "󰡨", 68},
	{"docx", "󰈬", 25},
	{"dot", "󱁉", 24},
	{"download", "", 79},
	{"drl", "", 217},
	{"dropbox", "", 26},
	{"dump", "", 188},
	{"dune", "", 94},
	{"dwg", "󰻫", 101},
	{"dwm", "", 31},
	{"dxf", "󰻫", 101},
	{"ebook", "", 180},
	{"ebuild", "", 60},
	{"edn", "", 73},
	{"eex", "", 139},
	{"ejs", "", 143},
	{"el", "", 103},
	{"elc", "", 103},
	{"elementary", "", 67},
	{"elf", "", 124},
	{"elm", "", 73},
	{"eln", "", 103},
	{"endeavour", "", 97},
	{"enlightenment", "", 231},
	{"env", "", 185},
	{"eot", "", 254},
	{"epp", "", 214},
	{"epub", "", 180},
	{"erb", "", 88},
	{"erl", "", 132},
	{"eslint.config.cjs", "", 55},
	{"eslint.config.js", "", 55},
	{"eslint.config.mjs", "", 55},
	{"eslint.config.ts", "", 55},
	{"ex", "", 139},
	{"exe", "", 124},
	{"exs", "", 139},
	{"ext_typoscript_setup.txt", "", 208},
	{"f#", "", 73},
	{"f3d", "󰻫", 101},
	{"f90", "󱈚", 96},
	{"favicon.ico", "", 143},
	{"fbx", "󰆧", 244},
	{"fcbak", "", 131},
	{"fcmacro", "", 131},
	{"fcmat", "", 131},
	{"fcparam", "", 131},
	{"fcscript", "", 131},
	{"fcstd", "", 131},
	{"fcstd1", "", 131},
	{"fctb", "", 131},
	{"fctl", "", 131},
	{"fdmdownload", "", 79},
	{"feature", "", 34},
	{"fedora", "", 17},
	{"fish", "", 59},
	{"flac", "", 31},
	{"flc", "", 254},
	{"flf", "", 254},
	{"fluxbox", "", 239},
	{"fnl", "", 224},
	{"fodg", "", 221},
	{"fodp", "", 179},
	{"fods", "", 113},
	{"fodt", "", 38},
	{"fp-info-cache", "", 231},
	{"fp-lib-table", "", 231},
	{"frag", "", 67},
	{"freebsd", "", 124},
	{"freecad.conf", "", 131},
	{"fs", "", 73},
	{"fsi", "", 73},
	{"fsscript", "", 73},
	{"fsx", "", 73},
	{"garuda", "", 32},
	{"gcode", "󰐫", 31},
	{"gd", "", 102},
	{"gemfile", "", 88},
	{"gemspec", "", 88},
	{"gentoo", "󰣨", 146},
	{"geom", "", 67},
	{"gif", "", 139},
	{"git", "", 166},
	{"glb", "", 215},
	{"gleam", "", 218},
	{"glsl", "", 67},
	{"gnome", "", 231},
	{"gnumakefile", "", 102},
	{"go", "", 38},
	{"go.mod", "", 38},
	{"go.sum", "", 38},
	{"go.work", "", 38},
	{"godot", "", 102},
	{"gpr", "", 102},
	{"gql", "", 169},
	{"gradle", "", 24},
	{"gradle-wrapper.properties", "", 24},
	{"gradle.properties", "", 24},
	{"gradlew", "", 24},
	{"graphql", "", 169},
	{"gresource", "", 231},
	{"groovy", "", 66},
	{"gruntfile.babel.js", "", 173},
	{"gruntfile.coffee", "", 173},
	{"gruntfile.js", "", 173},
	{"gruntfile.ts", "", 173},
	{"gtkrc", "", 231},
	{"guix", "", 220},
	{"gulpfile.babel.js", "", 167},
	{"gulpfile.coffee", "", 167},
	{"gulpfile.js", "", 167},
	{"gulpfile.ts", "", 167},
	{"gv", "󱁉", 24},
	{"gz", "", 178},
	{"h", "", 139},
	{"haml", "", 188},
	{"hbs", "", 172},
	{"heex", "", 139},
	{"hex", "", 27},
	{"hh", "", 139},
	{"hpp", "", 139},
	{"hrl", "", 132},
	{"hs", "", 139},
	{"htm", "", 166},
	{"html", "", 166},
	{"http", "", 31},
	{"huff", "󰡘", 61},
	{"hurl", "", 198},
	{"hx", "", 172},
	{"hxx", "", 139},
	{"hyperbola", "", 250},
	{"hypridle.conf", "", 37},
	{"hyprland", "", 37},
	{"hyprland.conf", "", 37},
	{"hyprlandd.conf", "", 37},
	{"hyprlock.conf", "", 37},
	{"hyprpaper.conf", "", 37},
	{"hyprsunset.conf", "", 37},
	{"i18n.config.js", "󰗊", 103},
	{"i18n.config.ts", "󰗊", 103},
	{"i3", "", 188},
	{"i3blocks.conf", "", 188},
	{"i3status.conf", "", 188},
	{"ical", "", 18},
	{"icalendar", "", 18},
	{"ico", "", 143},
	{"ics", "", 18},
	{"ifb", "", 18},
	{"ifc", "󰻫", 101},
	{"ige", "󰻫", 101},
	{"iges", "󰻫", 101},
	{"igs", "󰻫", 101},
	{"illumos", "", 202},
	{"image", "", 181},
	{"img", "", 181},
	{"import", "", 254},
	{"index.theme", "", 36},
	{"info", "", 230},
	{"ini", "", 102},
	{"ino", "", 73},
	{"ionic.config.json", "", 68},
	{"ipynb", "", 172},
	{"iso", "", 181},
	{"ixx", "", 73},
	{"jar", "", 216},
	{"java", "", 167},
	{"jenkinsfile", "", 167},
	{"jl", "", 139},
	{"jpeg", "", 139},
	{"jpg", "", 139},
	{"js", "", 143},
	{"json", "", 143},
	{"json5", "", 143},
	{"jsonc", "", 143},
	{"jsonl", "", 143},
	{"jsx", "", 38},
	{"justfile", "", 102},
	{"jwm", "", 32},
	{"jwmrc", "", 32},
	{"jxl", "", 139},
	{"kalgebrarc", "", 38},
	{"kali", "", 33},
	{"kbx", "󰯄", 102},
	{"kdb", "", 71},
	{"kdbx", "", 71},
	{"kdeglobals", "", 38},
	{"kdeneon", "", 37},
	{"kdenlive", "", 110},
	{"kdenlive-layoutsrc", "", 110},
	{"kdenliverc", "", 110},
	{"kdenlivetitle", "", 110},
	{"kicad_dru", "", 231},
	{"kicad_mod", "", 231},
	{"kicad_pcb", "", 231},
	{"kicad_prl", "", 231},
	{"kicad_pro", "", 231},
	{"kicad_sch", "", 231},
	{"kicad_sym", "", 231},
	{"kicad_wks", "", 231},
	{"ko", "", 188},
	{"kpp", "", 170},
	{"kra", "", 170},
	{"kritadisplayrc", "", 170},
	{"kritarc", "", 170},
	{"krz", "", 170},
	{"ksh", "", 59},
	{"kt", "", 99},
	{"kts", "", 99},
	{"kubuntu", "", 31},
	{"lck", "", 249},
	{"leap", "", 179},
	{"leex", "", 139},
	{"less", "", 60},
	{"lff", "", 254},
	{"lhs", "", 139},
	{"lib", "", 52},
	{"license", "", 179},
	{"license.md", "", 179},
	{"linux", "", 188},
	{"liquid", "", 107},
	{"lock", "", 249},
	{"locos", "", 178},
	{"log", "󰌱", 253},
	{"lrc", "󰨖", 214},
	{"lua", "", 74},
	{"luac", "", 74},
	{"luau", "", 39},
	{"lxde", "", 247},
	{"lxde-rc.xml", "", 245},
	{"lxle", "", 238},
	{"lxqt", "", 32},
	{"lxqt.conf", "", 32},
	{"m", "", 75},
	{"m3u", "󰲹", 175},
	{"m3u8", "󰲹", 175},
	{"m4a", "", 39},
	{"m4v", "", 172},
	{"mageia", "", 32},
	{"magnet", "", 124},
	{"makefile", "", 102},
	{"manjaro", "", 71},
	{"markdown", "", 253},
	{"mate", "", 149},
	{"material", "", 132},
	{"md", "", 253},
	{"md5", "󰕥", 103},
	{"mdx", "", 73},
	{"mint", "󰌪", 108},
	{"mix.lock", "", 139},
	{"mjs", "", 185},
	{"mk", "", 102},
	{"mkv", "", 172},
	{"ml", "", 173},
	{"mli", "", 173},
	{"mm", "", 73},
	{"mo", "", 104},
	{"mobi", "", 180},
	{"mojo", "", 202},
	{"mov", "", 172},
	{"mp3", "", 39},
	{"mp4", "", 172},
	{"mpp", "", 73},
	{"mpv.conf", "", 53},
	{"msf", "", 32},
	{"mts", "", 73},
	{"mustache", "", 173},
	{"mxlinux", "", 231},
	{"next.config.cjs", "", 231},
	{"next.config.js", "", 231},
	{"next.config.ts", "", 231},
	{"nfo", "", 230},
	{"nim", "", 184},
	{"nix", "", 110},
	{"nixos", "", 110},
	{"nobara", "", 231},
	{"node_modules", "", 161},
	{"norg", "", 67},
	{"nswag", "", 112},
	{"nu", "", 72},
	{"nuxt.config.cjs", "󱄆", 36},
	{"nuxt.config.js", "󱄆", 36},
	{"nuxt.config.mjs", "󱄆", 36},
	{"nuxt.config.ts", "󱄆", 36},
	{"o", "", 124},
	{"obj", "󰆧", 244},
	{"odf", "", 204},
	{"odg", "", 221},
	{"odin", "󰟢", 68},
	{"odp", "", 179},
	{"ods", "", 113},
	{"odt", "", 38},
	{"oga", "", 31},
	{"ogg", "", 31},
	{"ogv", "", 172},
	{"ogx", "", 172},
	{"openbsd", "", 178},
	{"opensuse", "", 106},
	{"opus", "", 31},
	{"org", "", 109},
	{"otf", "", 254},
	{"out", "", 124},
	{"package-lock.json", "", 88},
	{"package.json", "", 161},
	{"parabola", "", 103},
	{"parrot", "", 81},
	{"part", "", 79},
	{"patch", "", 59},
	{"pck", "", 102},
	{"pcm", "", 31},
	{"pdf", "", 124},
	{"php", "", 139},
	{"pkgbuild", "", 32},
	{"pl", "", 73},
	{"plasma", "", 32},
	{"platformio.ini", "", 172},
	{"playwright.config.cjs", "", 35},
	{"playwright.config.cts", "", 35},
	{"playwright.config.js", "", 35},
	{"playwright.config.mjs", "", 35},
	{"playwright.config.mts", "", 35},
	{"playwright.config.ts", "", 35},
	{"pls", "󰲹", 175},
	{"ply", "󰆧", 244},
	{"pm", "", 73},
	{"png", "", 139},
	{"pnpm-lock.yaml", "", 178},
	{"pnpm-workspace.yaml", "", 178},
	{"po", "", 31},
	{"pom.xml", "", 88},
	{"pop_os", "", 73},
	{"postmarketos", "", 34},
	{"pot", "", 31},
	{"pp", "", 214},
	{"ppt", "󰈧", 130},
	{"pptx", "󰈧", 130},
	{"prefab", "", 116},
	{"prettier.config.cjs", "", 68},
	{"prettier.config.js", "", 68},
	{"prettier.config.mjs", "", 68},
	{"prettier.config.ts", "", 68},
	{"prisma", "", 68},
	{"prisma.config.mts", "", 68},
	{"prisma.config.ts", "", 68},
	{"pro", "", 179},
	{"procfile", "", 139},
	{"prusaslicer.ini", "", 172},
	{"prusaslicergcodeviewer.ini", "", 172},
	{"ps1", "󰨊", 67},
	{"psb", "", 73},
	{"psd", "", 73},
	{"psd1", "󰨊", 103},
	{"psm1", "󰨊", 103},
	{"pub", "󰷖", 180},
	{"puppylinux", "", 145},
	{"pxd", "", 74},
	{"pxi", "", 74},
	{"py", "", 214},
	{"py.typed", "", 214},
	{"pyc", "", 222},
	{"pyd", "", 222},
	{"pyi", "", 214},
	{"pyo", "", 222},
	{"pyw", "", 74},
	{"pyx", "", 74},
	{"qm", "", 31},
	{"qml", "", 77},
	{"qrc", "", 77},
	{"qss", "", 77},
	{"qtile", "", 231},
	{"qtproject.conf", "", 77},
	{"qubesos", "", 68},
	{"query", "", 107},
	{"r", "󰟔", 31},
	{"rake", "", 88},
	{"rakefile", "", 88},
	{"rar", "", 178},
	{"rasi", "", 143},
	{"raspberry_pi", "", 125},
	{"razor", "󱦘", 56},
	{"rb", "", 88},
	{"readme", "󰂺", 254},
	{"readme.md", "󰂺", 254},
	{"redhat", "󱄛", 160},
	{"res", "", 167},
	{"resi", "", 168},
	{"river", "", 16},
	{"rkt", "󰘧", 124},
	{"rlib", "", 180},
	{"rmd", "", 73},
	{"robots.txt", "󰚩", 66},
	{"rocky", "", 36},
	{"rproj", "󰗆", 65},
	{"rs", "", 180},
	{"rss", "", 179},
	{"s", "", 31},
	{"sabayon", "", 251},
	{"sass", "", 168},
	{"sbt", "", 167},
	{"sc", "", 167},
	{"scad", "", 184},
	{"scala", "", 167},
	{"scm", "󰘧", 255},
	{"scss", "", 168},
	{"security", "󰒃", 145},
	{"security.md", "󰒃", 145},
	{"settings.gradle", "", 24},
	{"sh", "", 59},
	{"sha1", "󰕥", 103},
	{"sha224", "󰕥", 103},
	{"sha256", "󰕥", 103},
	{"sha384", "󰕥", 103},
	{"sha512", "󰕥", 103},
	{"sig", "󰘧", 173},
	{"signature", "󰘧", 173},
	{"skp", "󰻫", 101},
	{"slackware", "", 61},
	{"sldasm", "󰻫", 101},
	{"sldprt", "󰻫", 101},
	{"slim", "", 166},
	{"sln", "", 97},
	{"slnx", "", 97},
	{"slvs", "󰻫", 101},
	{"sml", "󰘧", 173},
	{"so", "", 188},
	{"sol", "", 73},
	{"solus", "", 59},
	{"spec.js", "", 143},
	{"spec.jsx", "", 38},
	{"spec.ts", "", 73},
	{"spec.tsx", "", 25},
	{"spx", "", 31},
	{"sql", "", 188},
	{"sqlite", "", 188},
	{"sqlite3", "", 188},
	{"srt", "󰨖", 214},
	{"ssa", "󰨖", 214},
	{"ste", "󰻫", 101},
	{"step", "󰻫", 101},
	{"stl", "󰆧", 244},
	{"stories.js", "", 204},
	{"stories.jsx", "", 204},
	{"stories.mjs", "", 204},
	{"stories.svelte", "", 204},
	{"stories.ts", "", 204},
	{"stories.tsx", "", 204},
	{"stories.vue", "", 204},
	{"stp", "󰻫", 101},
	{"strings", "", 31},
	{"styl", "", 107},
	{"sub", "󰨖", 214},
	{"sublime", "", 173},
	{"suo", "", 97},
	{"sv", "󰍛", 29},
	{"svelte", "", 202},
	{"svelte.config.js", "", 202},
	{"svg", "󰜡", 215},
	{"svgz", "󰜡", 215},
	{"svh", "󰍛", 29},
	{"svx", "", 197},
	{"sway", "", 100},
	{"swift", "", 173},
	{"sxhkdrc", "", 235},
	{"sym-lib-table", "", 231},
	{"t", "", 73},
	{"tails", "", 60},
	{"tailwind.config.js", "󱏿", 38},
	{"tailwind.config.mjs", "󱏿", 38},
	{"tailwind.config.ts", "󱏿", 38},
	{"tbc", "󰛓", 25},
	{"tcl", "󰛓", 25},
	{"templ", "", 178},
	{"terminal", "", 35},
	{"test.js", "", 143},
	{"test.jsx", "", 38},
	{"test.ts", "", 73},
	{"test.tsx", "", 25},
	{"tex", "", 58},
	{"tf", "", 62},
	{"tfvars", "", 62},
	{"tgz", "", 178},
	{"tmpl", "", 178},
	{"tmux", "", 34},
	{"tmux.conf", "", 34},
	{"tmux.conf.local", "", 34},
	{"toml", "", 130},
	{"torrent", "", 79},
	{"tres", "", 102},
	{"trisquel", "", 25},
	{"ts", "", 73},
	{"tscn", "", 102},
	{"tsconfig", "", 208},
	{"tsconfig.json", "", 73},
	{"tsx", "", 25},
	{"ttf", "", 254},
	{"tumbleweed", "", 73},
	{"twig", "", 107},
	{"txt", "󰈙", 113},
	{"txz", "", 178},
	{"typ", "", 37},
	{"typoscript", "", 208},
	{"ubuntu", "", 166},
	{"ui", "", 26},
	{"unity", "", 231},
	{"unlicense", "", 179},
	{"v", "󰍛", 29},
	{"vagrantfile", "", 27},
	{"vala", "", 97},
	{"vanillaos", "", 179},
	{"vercel.json", "", 231},
	{"vert", "", 67},
	{"vh", "󰍛", 29},
	{"vhd", "󰍛", 29},
	{"vhdl", "󰍛", 29},
	{"vi", "", 178},
	{"vim", "", 29},
	{"vite.config.cjs", "", 214},
	{"vite.config.cts", "", 214},
	{"vite.config.js", "", 214},
	{"vite.config.mjs", "", 214},
	{"vite.config.mts", "", 214},
	{"vite.config.ts", "", 214},
	{"vitest.config.cjs", "", 106},
	{"vitest.config.cts", "", 106},
	{"vitest.config.js", "", 106},
	{"vitest.config.mjs", "", 106},
	{"vitest.config.mts", "", 106},
	{"vitest.config.ts", "", 106},
	{"vlcrc", "󰕼", 172},
	{"void", "", 23},
	{"vsh", "", 67},
	{"vsix", "", 97},
	{"vue", "", 107},
	{"wasm", "", 62},
	{"wav", "", 39},
	{"webm", "", 172},
	{"webmanifest", "", 185},
	{"webp", "", 139},
	{"webpack", "󰜫", 73},
	{"weston.ini", "", 214},
	{"windows", "", 38},
	{"wma", "", 39},
	{"wmv", "", 172},
	{"woff", "", 254},
	{"woff2", "", 254},
	{"workspace", "", 113},
	{"wrangler.jsonc", "", 172},
	{"wrangler.toml", "", 172},
	{"wrl", "󰆧", 244},
	{"wrz", "󰆧", 244},
	{"wv", "", 39},
	{"wvc", "", 39},
	{"x", "", 75},
	{"xaml", "󰙳", 56},
	{"xcf", "", 59},
	{"xcplayground", "", 173},
	{"xcstrings", "", 31},
	{"xdph.conf", "", 37},
	{"xerolinux", "", 104},
	{"xfce", "", 38},
	{"xls", "󰈛", 29},
	{"xlsx", "󰈛", 29},
	{"xm", "", 73},
	{"xml", "󰗀", 173},
	{"xmobarrc", "", 167},
	{"xmobarrc.hs", "", 167},
	{"xmonad", "", 167},
	{"xmonad.hs", "", 167},
	{"xorg.conf", "", 166},
	{"xpi", "", 196},
	{"xsettingsd.conf", "", 166},
	{"xslt", "󰗀", 74},
	{"xul", "", 173},
	{"xz", "", 178},
	{"yaml", "", 160},
	{"yml", "", 160},
	{"zig", "", 178},
	{"zip", "", 178},
	{"zorin", "", 38},
	{"zsh", "", 113},
	{"zst", "", 178},
	{"🔥", "", 202},
}

var (
	lastExt   string
	lastMatch IconEntry
	hasCache  bool
)

func compareStringBytes(s string, b []byte) int {
	minLen := len(s)
	if len(b) < minLen {
		minLen = len(b)
	}
	for i := 0; i < minLen; i++ {
		if s[i] < b[i] {
			return -1
		} else if s[i] > b[i] {
			return 1
		}
	}
	if len(s) < len(b) {
		return -1
	} else if len(s) > len(b) {
		return 1
	}
	return 0
}

func getIcon(filename []byte) IconEntry {
	// 1. Try exact filename match
	idx, found := slices.BinarySearchFunc(FilenameIcons, filename, func(entry IconEntry, target []byte) int {
		return compareStringBytes(entry.Key, target)
	})
	if found {
		return FilenameIcons[idx]
	}

	// 2. Extract extension and try extension match
	dotIdx := -1
	for i := len(filename) - 1; i >= 0; i-- {
		if filename[i] == '.' {
			dotIdx = i
			break
		}
	}
	if dotIdx >= 0 && dotIdx < len(filename)-1 {
		ext := filename[dotIdx+1:]
		
		// 1-element cache check
		if hasCache && string(ext) == lastExt {
			return lastMatch
		}

		idx, found = slices.BinarySearchFunc(ExtensionIcons, ext, func(entry IconEntry, target []byte) int {
			return compareStringBytes(entry.Key, target)
		})
		
		var match IconEntry
		if found {
			match = ExtensionIcons[idx]
		} else {
			match = IconEntry{Key: "", Icon: "", Color: 231}
		}

		// Cache
		lastExt = string(ext)
		lastMatch = match
		hasCache = true

		return match
	}

	return IconEntry{Key: "", Icon: "", Color: 231}
}

func formatLineWithIcon(data []byte) []byte {
	// Find last '/' to split filename and dir
	lastSlash := -1
	for i := len(data) - 1; i >= 0; i-- {
		if data[i] == '/' || data[i] == '\\' {
			lastSlash = i
			break
		}
	}

	var filename []byte
	var dir []byte
	if lastSlash >= 0 {
		filename = data[lastSlash+1:]
		dir = data[:lastSlash]
	} else {
		filename = data
	}

	match := getIcon(filename)

	var buf bytes.Buffer
	buf.Grow(45 + len(filename) + len(dir))

	buf.WriteString("\x1b[38;5;")
	buf.WriteString(strconv.Itoa(match.Color))
	buf.WriteByte('m')
	buf.WriteString(match.Icon)
	buf.WriteString("\x1b[0m  ")
	buf.Write(filename)
	buf.WriteString("  ")

	if len(dir) > 0 {
		buf.WriteString("\x1b[38;5;244;3m")
		buf.Write(dir)
		buf.WriteString("\x1b[0m")
	}

	return buf.Bytes()
}

func pushMruFiles(cl *ChunkList) {
	mruEnv := os.Getenv("FZF_MRU_LIST")
	if mruEnv == "" {
		return
	}
	cwd, err := os.Getwd()
	parts := strings.Split(mruEnv, ";")
	for _, file := range parts {
		if file == "" {
			continue
		}
		relFile := file
		if err == nil && strings.HasPrefix(file, cwd) {
			if len(file) > len(cwd) {
				relFile = file[len(cwd)+1:]
			}
		}
		
		formatted := formatLineWithIcon([]byte(relFile))
		cl.Push(formatted)
	}
}

func ExtractPathFromFormatted(data []byte) string {
	if len(data) == 0 {
		return ""
	}
	if data[0] != '\x1b' {
		return string(data)
	}

	cleaned := stripAnsi(data)
	parts := strings.Split(cleaned, "  ")
	if len(parts) < 2 {
		return cleaned
	}

	filename := strings.TrimSpace(parts[1])
	if len(parts) >= 3 {
		dir := strings.TrimSpace(parts[2])
		if dir != "" {
			return dir + "/" + filename
		}
	}
	return filename
}

func stripAnsi(data []byte) string {
	var buf strings.Builder
	inEscape := false
	for i := 0; i < len(data); i++ {
		if data[i] == '\x1b' {
			inEscape = true
			continue
		}
		if inEscape {
			if data[i] == 'm' {
				inEscape = false
			}
			continue
		}
		buf.WriteByte(data[i])
	}
	return buf.String()
}

func isMruFile(data []byte) bool {
	if len(algo.MruMap) == 0 {
		return false
	}
	path := ExtractPathFromFormatted(data)
	_, ok := algo.GetMruRank(path)
	return ok
}

func formatLineWithDummyIcon(data []byte) []byte {
	if len(data) > 0 && data[0] == '\x1b' {
		return data
	}
	// Find last '/' to split filename and dir
	lastSlash := -1
	for i := len(data) - 1; i >= 0; i-- {
		if data[i] == '/' || data[i] == '\\' {
			lastSlash = i
			break
		}
	}

	var filename []byte
	var dir []byte
	if lastSlash >= 0 {
		filename = data[lastSlash+1:]
		dir = data[:lastSlash]
	} else {
		filename = data
	}

	var buf bytes.Buffer
	buf.Grow(20 + len(filename) + len(dir))

	buf.WriteString("\x1b[38;5;231m\x1b[0m  ")
	buf.Write(filename)
	buf.WriteString("  ")

	if len(dir) > 0 {
		buf.WriteString("\x1b[38;5;244;3m")
		buf.Write(dir)
		buf.WriteString("\x1b[0m")
	}

	return buf.Bytes()
}

func getRealIconFromRunes(runes []rune) IconEntry {
	if len(runes) < 4 || runes[0] != '' {
		return IconEntry{Key: "", Icon: "", Color: 231}
	}

	// Scan until double space
	endIdx := -1
	for i := 3; i < len(runes)-1; i++ {
		if runes[i] == ' ' && runes[i+1] == ' ' {
			endIdx = i
			break
		}
	}
	if endIdx < 0 {
		endIdx = len(runes)
	}

	filenameRunes := runes[3:endIdx]
	filenameStr := string(filenameRunes)
	return getIcon([]byte(filenameStr))
}
