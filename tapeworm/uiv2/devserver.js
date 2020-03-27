var bs = require('browser-sync').create();

process.stdin.on('data', function (data) {
    console.log(`Received notification: ${data.toString('utf8').replace("\n", "")}`);

    bs.reload('*')
});

bs.init({
    server: 'tapeworm/uiv2/',
    open: false,
})
