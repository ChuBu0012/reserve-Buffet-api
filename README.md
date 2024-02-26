
``` 
websocket -> /ws
updatestate -> /updatestate/:id POST Method
{
    "phone": "099-999-9999",
    "status": (0,1,2), 0 = ว่าง, 1 = จองเเล้ว, 2 = โต๊ะไม่ว่าง
    "startTime": "2024-2-19 16.00",
    "endTime": "2024-2-19 18.50"
}
gettables -> /gettables GET Method
```

``` jsx
import { useEffect, useState } from 'react';

function App() {
  const [wsData, setWsData] = useState([]);
  const [ws, setWs] = useState(null);

  useEffect(() => {
    const fetchData = async () => {
      const response = await fetch('http://localhost:8080/gettables');
      const data = await response.json();
      setWsData(data)
      console.log(data)
    };

    fetchData();

    const ws = new WebSocket('ws://localhost:8080/ws');

    ws.onopen = () => {
      console.log('WebSocket is connected.');
    };

    ws.onmessage = (event) => {
      console.log(JSON.parse(event.data));
      const newData = JSON.parse(event.data);
      setWsData(oldData => oldData.map(ele => ele.tableid === newData.tableid ? newData : ele));
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


  return (
    <>
    <table cellSpacing={10} >
      <thead></thead>
      <tbody>
      {
        wsData.map((ele)=>{
          return <tr key={ele.tableid}>
            <td>{ele.tableid}</td>
            <td>{ele.phone}</td>
            <td>{ele.status}</td>
            <td>{ele.code}</td>
            <td>{ele.startTime}</td>
            <td>{ele.endTime}</td>
            </tr>
        })
      }

      </tbody>
      <tfoot></tfoot>
    </table>
    </>
  );
}

export default App;

```