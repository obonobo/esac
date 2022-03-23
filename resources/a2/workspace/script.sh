main() (
    strings=(
        "<START>"
        "<assignStat>"
        "<implDef>"
        "<returnType>"
        "<term>"
        "<rept-variable0>"
        "<idnest>"
        "<memberDecl>"
        "<rept-implDef3>"
        "<type>"
        "<rept-variable2>"
        "<arraySize>"
        "<variable>"
        "<sign>"
        "<assignOp>"
        "<rightrec-arithExpr>"
        "<addOp>"
        "<rept-idnest1>"
        "<funcDef>"
        "<funcBody>"
        "<rept-fParamsTail4>"
        "<rept-opt-structDecl22>"
        "<funcHead>"
        "<statBlock>"
        "<arithExpr>"
        "<varDecl>"
        "<functionCall>"
        "<rept-fParams3>"
        "<aParamsTail>"
        "<varDeclOrStat>"
        "<multOp>"
        "<relOp>"
        "<opt-structDecl2>"
        "<rept-structDecl4>"
        "<fParams>"
        "<indice>"
        "<rept-fParams4>"
        "<rept-aParams1>"
        "<structDecl>"
        "<visibility>"
        "<relExpr>"
        "<rept-functionCall0>"
        "<aParams>"
        "<expr>"
        "<rightrec-term>"
        "<prog>"
        "<rept-prog0>"
        "<rept-funcBody1>"
        "<funcDecl>"
        "<rept-varDecl4>"
        "<statement>"
        "<factor>"
        "<fParamsTail>"
        "<structOrImplOrFunc>"
        "<rept-statBlock1>"
    )

    for s in "${strings[@]}"; do
        ss="${s#<}"
        ss="${ss%>}"
        ss="$(echo $ss | tr [a-z] [A-Z] | tr -d '-')"
        echo "$ss Kind = \"$s\""
    done
)

main
# [[ ${BASH_SOURCE[0]} == $0 ]] && main "$@"