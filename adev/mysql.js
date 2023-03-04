import sql from 'k6/x/sql';

// The second argument is a MySQL connection string, e.g.
// myuser:mypass@tcp(127.0.0.1:3306)/mydb
const db = sql.open('mysql', 'root:root@tcp(127.0.0.1:3306)/bbgo');


export default function () {
    // //  db.exec("INSERT INTO keyvalues (key, value) VALUES('plugin-name', 'k6-plugin-sql');");
    // //
    // // let results = sql.query(db, "SELECT * FROM orders limit 10;");
    // // for (const row of results) {
    // //     // Convert array of ASCII integers into strings. See https://github.com/grafana/xk6-sql/issues/12
    // //     console.log(`key: ${String.fromCharCode(...row.gid)}, value: ${String.fromCharCode(...row.order_id)}`);
    // // }
    db.exec("INSERT INTO keyvalues (`key`, value) VALUES('plugin-name', 'k6-plugin-sql');");

    let results = sql.query(db, "SELECT * FROM orders");
    for (const row of results) {
        // Convert array of ASCII integers into strings. See https://github.com/grafana/xk6-sql/issues/12
        console.log(`key: ${String.fromCharCode(...row.gid)}, value: ${String.fromCharCode(...row.order_id)}`);
    }
    return 333
}
