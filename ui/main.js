var ctx = document.getElementById('myChart').getContext('2d');
var chart = new Chart(ctx, {
    type: 'line',
    data: {
        datasets: [{
            label: 'GET',
            data: [],
            lineTension: 0.1,
            backgroundColor: 'rgba(85, 239, 196, 0.3)',
            borderColor: '#00b894',
            borderWidth: 1,
            spanGaps: false,
        }, {
            label: 'POST',
            data: [],
            lineTension: 0.1,
            backgroundColor: 'rgba(116, 185, 255,0.3)',
            borderColor: '#0984e3',
            borderWidth: 1,
            spanGaps: false,
        },
        {
            label: 'PUT',
            data: [],
            lineTension: 0.1,
            backgroundColor: 'rgba(162, 155, 254,0.3)',
            borderColor: '#6c5ce7',
            borderWidth: 1,
            spanGaps: false,
        },
        {
            label: 'DELETE',
            data: [],
            lineTension: 0.1,
            backgroundColor: 'rgba(255, 118, 117, 0.3)',
            borderColor: '#d63031',
            borderWidth: 1,
            spanGaps: false,
        }]
    },
    options: {
        scales: {
            xAxes: [{
                type: 'time',
                distribution: 'series',
                time: {
                    tooltipFormat: 'HH:mm',
                    unit: 'minute',
                    displayFormats: {
                        'minute': 'HH:mm',
                        'hour': 'HH:mm'
                    },
                }
            }],
            yAxes: [{
                ticks: {
                    beginAtZero: true
                },
            }]
        }
    }
});


const socket = io("http://localhost:3000");

socket.on("connect", () => {
    console.log("Connected to the websocket.");
});

socket.on("data", data => {
    if (data.length) {
        const ONE_HOUR = 60 * 60 * 1000;
        const serie = [];
        const oneHourLater = Date.now() - ONE_HOUR;

        for (let index = 0; index < 60; index++) {
            const dMinuteAdded = new Date(oneHourLater)
            dMinuteAdded.setMinutes(dMinuteAdded.getMinutes() + index)

            serie.push({
                t: dMinuteAdded.toLocaleString(),
                y: 0
            })
        }

        const dr = {
            GET: [...JSON.parse(JSON.stringify(serie))],
            POST: [...JSON.parse(JSON.stringify(serie))],
            PUT: [...JSON.parse(JSON.stringify(serie))],
            DELETE: [...JSON.parse(JSON.stringify(serie))],
        }

        data.map(d => {
            const t = new Date(d.timestamp * 1000)

            for (let i = 1; i < dr[d.type].length; i++) {
                if (new Date(dr[d.type][i].t) > t) {
                    dr[d.type][i - 1].y = dr[d.type][i - 1].y + parseInt(d.spent)
                    break;
                }
            }

        })

        chart.data.datasets.forEach((dataset) => {
            dataset.data.push(...dr[dataset.label]);
        });

        chart.update();
    }
})

function spanData() {

    chart.data.datasets.forEach((dataset) => {
        while (dataset.data.length > 60) {
            dataset.data.shift();
        }
    });

}

socket.on("new", d => {

    const t = new Date(d.timestamp * 1000)

    for (let j = 0; j < chart.data.datasets.length; j++) {
        const dataset = chart.data.datasets[j];
        if (dataset.label === d.type) {

            let flag = true
            for (let index = dataset.data.length - 1; index > -1; index--) {
                const current = new Date(dataset.data[index].t);
                if (current.getMinutes() == t.getMinutes() && current <= t) {
                    dataset.data[index].y = dataset.data[index].y + parseInt(d.spent)
                    flag = false
                    break
                } else if (current < t) {
                    break
                }

            }

            if (flag) {

                chart.data.datasets.forEach((dataset) => {
                    dataset.data.push({
                        t: t.toISOString(),
                        y: dataset.label === d.type ? parseInt(d.spent) : 0
                    });
                });
            }

            break
        }
    }

    spanData()
    chart.update();
})


setInterval(()=>{
    const timeNow = new Date()    
    chart.data.datasets.forEach((dataset) => {
        if(new Date(dataset.data[dataset.data.length - 1].t).getMinutes() != timeNow.getMinutes()){
            dataset.data.push({
                t: timeNow.toLocaleString(),
                y: 0
            });
        }
    });

    spanData()
    chart.update()

}, 1000 * 60)