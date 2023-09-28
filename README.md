# exhy-cloud
An distributed orchestrator consists of basic components. 

## Worker
1. Run task as Docker containers.
2. Accept tasks to run from a manager.
3. Provide relevant statistics to manager for the purpose of scheduling tasks.
4. Keep track of its tasks and their state.

## Manager
1. Accept requests from users to start and stop tasks.
2. Schedule tasks onto worker machines.
3. Keep track of tasks. their states, and the machine on which they run.