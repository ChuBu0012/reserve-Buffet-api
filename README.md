
``` 
websocket -> :8080/ws
setname -> :8080/setname POST Method
getname -> :8080/getname GET Method
```

``` jsx
  const [ws, setWs] = useState(null);

  useEffect(() => {
    const ws = new WebSocket('ws://localhost:8080/ws');

    ws.onopen = () => {
      console.log('WebSocket is connected.');
    };

    ws.onmessage = (event) => {
      console.log('Received data from server:', event.data);
      setWsData(event.data);
    };

    ws.onerror = (error) => {
      console.log('WebSocket error:', error);
    };

    ws.onclose = () => {
      console.log('WebSocket is closed.');
    };

    setWs(ws);
    
    return () => {
      ws.close();
    };
  }, []);
```