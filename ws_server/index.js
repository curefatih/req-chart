
const httpServer = require("http").createServer();
const io = require("socket.io")(httpServer, {
    cors: {
        origin: "*",
        methods: ["GET", "POST"]
    }
});

const { Client } = require('pg');

const client = new Client({
    user: 'postgres',
    host: 'db',
    database: 'postgres',
    password: 'postgres',
    port: 5432,
});

client.connect();

client.query("LISTEN new_data")

io.on("connection", (socket) => {
    console.log("socket connected", socket.id);
    const ONE_HOUR = 60 * 60 * 1000;
    const oneHourBefore = Math.ceil((new Date().getTime() - ONE_HOUR) / 1000);
    client.query('SELECT * FROM requests WHERE timestamp > ' + oneHourBefore, (err, res) => {
        if (err) {
            console.log(err.stack)
        } else {
            socket.emit("data", res.rows)
        }
    })

    client.on("notification", async (data) => {
        socket.emit("new", JSON.parse(data.payload))
    })
});

httpServer.listen(3000);