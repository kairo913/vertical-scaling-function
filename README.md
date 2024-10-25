# Function Resize Shape

Update instance shape to specified size

# Installation

1. Create Compartment

2. Create VCN (must use to the compartment created in step 1)

3. Create the Instance you want to resize
    - Shape must match either one you want to specify and must use to the Compartment and VCN created in step 1,2.

4. Create Function
    - Create Application (must use to the Compartment and VCN created in step 1,2)
    - Set the following environment variables
        ```
        INSTANCE_ID: <instance id>
        ACTIVE_OCPU: <OCPU amount>
        ACTIVE_MEMORY: <Memory amount (GiB)>
        INACTIVE_OCPU: <OCPU amount while inavtive mode>
        INACTIVE_MEMORY: <Memory amount (GiB) while inactive mode>
        ```
    - Follow start guide to setup function deploy
    - Open Cloud Shell and run below commands.

        ```
        git clone https://github.com/kairo913/  function-resize-shape.git
        cd function-resize-shape/fn-resize-shape
        fn -v deploy --app <application name>
        ```

5. Create Resource Schedule
    - Select "Start" as the action to execute, and specify the created Function as the resource.
    - If you want to switch the shape size of an instance multiple times, create multiple schedules.

6. Create Dynamic-Group
    - Use the following matching rules
        ```
        ALL {resource.type = 'fnfunc', resource.id = '<function ocid>'}
        ALL {resource.type='resourceschedule', resource.id='<resource schedule ocid>'}
        ```

    - Add the following rules for the number of additional schedules created
        ```
        ALL {resource.type='resourceschedule', resource.id='<resource schedule ocid>'}
        ```