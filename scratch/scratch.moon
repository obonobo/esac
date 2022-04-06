% Main:
	% PROG
	entry
	addi	r14, r0, topaddr
	
	% INTNUM 10
	addi	r12, r0, 10
	sw	t0(r0), r12
	
	% ASSIGN x = t0
	lw	r12, t0(r0)
	sw	x(r0), r12
	
	% ASSIGN y = x
	lw	r12, x(r0)
	sw	y(r0), r12
	
	% MULT
	lw	r11, y(r0)
	lw	r10, y(r0)
	mul	r12, r11, r10
	sw	t1(r0), r12
	
	% ASSIGN y = t1
	lw	r12, t1(r0)
	sw	y(r0), r12
	
	% PLUS
	lw	r11, x(r0)
	lw	r10, y(r0)
	add	r12, r11, r10
	sw	t2(r0), r12
	
	% WRITE(t2)
	lw	r12, t2(r0)
	sw	-8(r14), r12	% intstr arg1
	addi	r12, r0, wbuf
	sw	-12(r14), r12	% intstr arg2
	jl	r15, intstr	% Procedure call intstr
	sw	-8(r14), r13	% putstr arg1
	jl	r15, putstr	% Procedure call putstr
	hlt

% Data:
x	res	4		% Space for variable x
y	res	4		% Space for variable y
t0	res	4		% Space for variable t0
t1	res	4		% Space for variable t1
t2	res	4		% Space for variable t2
wbuf	res	32		% Buffer for printing

