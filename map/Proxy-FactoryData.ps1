while(1){ 
    wget -OutFile factory.json http://localhost:8090/getFactory; 
    Start-Sleep -Seconds 20; 
}