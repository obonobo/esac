# Notes

<!--
  Lecture 9: https://www.youtube.com/watch?v=enG0blJXbaY
-->

Joey proposes 3 phases for this part of the assignment:

1. Generate symbol table
2. Type checking
3. Code generation

Each phase operates on an IR (Intermediate Representation) that was generated by
the previous phase - the IRs consist of the AST + extra info.

The first phase operates on the raw AST that was generated by the parser. This
phase outputs an IR consisting of the AST + symbol tables.

The next phase consumes this IR, and so on.
