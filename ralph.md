PSEUDOCODE - this is your prompt - PSEUDOCODE


NOTE no help is coming for you, this is an automated system

SUB complete_task:
    format new and existing code
    DO:
        run the tests
        fix if needed
    WHILE the tests are broken;
    
    DO:
        build the project
    WHILE the build is failing;

    git add .
    git commit <conventional commit style>
    git push

    EXIT ALL, STOP WORK;
END SUB

SUB complete_simple_task:
    git add .
    git commit <conventional commit style>
    git push

    EXIT ALL, STOP WORK;
END SUB

SUB main
    read the in progress tickets
    if there are tickets in progress
        complete exactly one ticket 
        move to done
        CALL complete_task
    
    read the todo tickets
    if there are todo tickets
        move exactly one ticket to in progress
        complete that ticket
        move to done
        CALL complete_task
    
    read the in progress ERDs
    read the in progress PRDs
    if there are no todo tickets AND there is precisely one ERD in progress AND there is precisely one PRD in progress
        move that ERD to done
        move that PRD to done
        CALL complete_simple_task
    
    read the todo ERDs
    read the todo PRDs
    if there are no ERDs in progress AND there are no PRDs in progress AND there is at least one ERD/PRD pair in todo
        find the first ERD/PRD pair in todo (by precise order)
        move that ERD to in progress
        move that PRD to in progress
        CALL complete_simple_task
    
    

END

CALL main