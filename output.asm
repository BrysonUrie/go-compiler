section .data
fmtstr: db "%s", 10, 0
fmtint: db "%d", 10, 0
fmtfloat: db "%f", 10, 0
fmtintin: db "%d", 0
fmtfloatin: db "%f", 0
float1: dd 0.0
section .text
extern printf
global main
main:
push RBP
mov RBP, RSP
sub RSP, 64
pop RAX
mov dword [edp-4], RAX
mov RAX, 2
mov RCX, dword [rbp-4]
xor RDX, RDX
div RCX
push RAX
pop RAX
mov dword [edp-4], RAX
