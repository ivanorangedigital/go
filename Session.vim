let SessionLoad = 1
let s:so_save = &g:so | let s:siso_save = &g:siso | setg so=0 siso=0 | setl so=-1 siso=-1
let v:this_session=expand("<sfile>:p")
silent only
silent tabonly
cd ~/Desktop/webDev/orangebud/last/backend
if expand('%') == '' && !&modified && line('$') <= 1 && getline(1) == ''
  let s:wipebuf = bufnr('%')
endif
let s:shortmess_save = &shortmess
if &shortmess =~ 'A'
  set shortmess=aoOA
else
  set shortmess=aoO
endif
badd +7 makefile
badd +25 cmd/web/templates.go
badd +24 readme.md
badd +20 ui/html/pages/home.page.tmpl
badd +18 dev-ws-server/server.js
badd +24 cmd/web/helpers.go
badd +1 ui/html/partials/default-header.partial.tmpl
badd +37 ui/html/layouts/base.layout.tmpl
badd +2 ui/html/partials/default-footer.partial.tmpl
badd +19 cmd/web/middlewares.go
badd +1 ui/static/js/index.js
badd +6 cmd/web/main.go
badd +1 ui/static/js/dev-watch.js
badd +13 cmd/web/handlers.go
argglobal
%argdel
$argadd NvimTree_1
edit cmd/web/handlers.go
let s:save_splitbelow = &splitbelow
let s:save_splitright = &splitright
set splitbelow splitright
wincmd _ | wincmd |
vsplit
1wincmd h
wincmd w
let &splitbelow = s:save_splitbelow
let &splitright = s:save_splitright
wincmd t
let s:save_winminheight = &winminheight
let s:save_winminwidth = &winminwidth
set winminheight=0
set winheight=1
set winminwidth=0
set winwidth=1
exe 'vert 1resize ' . ((&columns * 30 + 85) / 170)
exe 'vert 2resize ' . ((&columns * 139 + 85) / 170)
argglobal
enew
file NvimTree_1
balt makefile
setlocal fdm=manual
setlocal fde=0
setlocal fmr={{{,}}}
setlocal fdi=#
setlocal fdl=0
setlocal fml=1
setlocal fdn=20
setlocal nofen
lcd ~/Desktop/webDev/orangebud/last/backend
wincmd w
argglobal
balt ~/Desktop/webDev/orangebud/last/backend/cmd/web/main.go
setlocal fdm=manual
setlocal fde=0
setlocal fmr={{{,}}}
setlocal fdi=#
setlocal fdl=0
setlocal fml=1
setlocal fdn=20
setlocal fen
silent! normal! zE
9,45fold
let &fdl = &fdl
let s:l = 9 - ((8 * winheight(0) + 26) / 53)
if s:l < 1 | let s:l = 1 | endif
keepjumps exe s:l
normal! zt
keepjumps 9
normal! 0
lcd ~/Desktop/webDev/orangebud/last/backend
wincmd w
2wincmd w
exe 'vert 1resize ' . ((&columns * 30 + 85) / 170)
exe 'vert 2resize ' . ((&columns * 139 + 85) / 170)
tabnext 1
if exists('s:wipebuf') && len(win_findbuf(s:wipebuf)) == 0 && getbufvar(s:wipebuf, '&buftype') isnot# 'terminal'
  silent exe 'bwipe ' . s:wipebuf
endif
unlet! s:wipebuf
set winheight=1 winwidth=20
let &shortmess = s:shortmess_save
let &winminheight = s:save_winminheight
let &winminwidth = s:save_winminwidth
let s:sx = expand("<sfile>:p:r")."x.vim"
if filereadable(s:sx)
  exe "source " . fnameescape(s:sx)
endif
let &g:so = s:so_save | let &g:siso = s:siso_save
set hlsearch
nohlsearch
doautoall SessionLoadPost
unlet SessionLoad
" vim: set ft=vim :
