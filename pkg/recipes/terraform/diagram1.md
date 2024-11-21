+-----------------------------------------------------+
|       Terraform (e *executor) Deploy()              |
|                                                     |
+-----------------------------------------------------+
                         |
                         v
+-----------------------------------------------------+     +------------------------------+     Move configuration /recipes/providers/main.tf
| Step 1: Create Temporary Data Loader Pod            |      | =>  |    Copy Data into      |     into Volume
| -                                                    |     |        Volume (diagram2)     |     => first iteration, tried using ConfigMap     
+-----------------------------------------------------+     +------------------------------+     ran into size constraints
                         |                                                                       * rbac permissions to be able to create PVC/PV  
                         v
+-----------------------------------------------------+
| Step 2: Create Terraform Job                        |
| - Generate a  job name using the environment        |     *Using (kubernetes secret) backend. Set up service account/roles/rolebindings  
|   recipe name                                       |     for access to secret from container spun off in Job. Set job to run under serviceacc. 
| - Call createTerraformJob                           |
|                                                     |
+-----------------------------------------------------+
                         |
                         v
+-----------------------------------------------------+
| Step 3: Wait for Job Completion                     |
| - Call waitForJobCompletion with the job name       |
|             |
+-----------------------------------------------------+
                         |
                         v
+-----------------------------------------------------+
| Step 5: Retrieve Terraform Outputs                  |
| - Call retrieveTerraformOutputs                     |
| - Handle errors during output retrieval             |
+-----------------------------------------------------+
                         |
                         v
+-----------------------------------------------------+
| Step 6: Return Terraform State                      |
| - Return the retrieved Terraform state object       |
| - End of Apply function                             |
+-----------------------------------------------------+