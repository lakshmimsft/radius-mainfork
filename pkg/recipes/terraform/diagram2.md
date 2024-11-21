+-----------------------------------------------------+
|         Function: copyDataToPVC(ctx, options)       |
+-----------------------------------------------------+
                         |
                         v
+-----------------------------------------------------+
| Step 1: Create Exec Command                         |
| - Define the command to extract data in the pod     |
|   (e.g., tar extraction command)                    |
+-----------------------------------------------------+
                         |
                         v
                         
                         v
+-----------------------------------------------------+
| Step 3: Create Tar Stream from Source Path          |
| - Call createTarStream with the source path         |
| - Generates an in-memory tar archive of the files   |
+-----------------------------------------------------+
                         |
                         v
+-----------------------------------------------------+
| Step 4: Stream Tar Archive into Pod                 |
| - Stream the tar archive to the pod's stdin         |
| - The pod extracts the files into the target path   |
| - Handle any execution errors                       |
+-----------------------------------------------------+