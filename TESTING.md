## How to run pre-written tests
<strong>BEFORE YOU RUN:</strong></br>
Currently, if one of the test url's are already stored in the DB, then it will return as a FAIL,
when it just means that it is not going to overwrite the data because it already exists. If you want to ensure 100% your
tests are correct, use ```mv``` to move the current ```chroma``` directory outside of the Silvercord directory,
then run your tests. After passing, you can replace the new directory generated with the "moved" one after.
1. Make sure Go is installed </br>
2. Run ```go mod init <module-name-here>```</br>
Example: ```go mod init db_test``` </br>
3. Run ```go test```</br>
4. Clean up Go environment: ```rm go.mod```
