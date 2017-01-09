/**
 * Created by dafengge0913 on 2017/1/9.
 */

function sendName() {

    var name = document.getElementById("input_name").value;
    //console.log(name);
    var xhr = new XMLHttpRequest();
    xhr.open('post', '/helloAjax', true);
    xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
    var req = {};
    req.name = name;
    xhr.send(JSON.stringify(req));
    xhr.onreadystatechange = function () {
        if (xhr.readyState == 4 && xhr.status == 200) {
            document.getElementById('result').innerHTML = xhr.responseText;
        }
    }
}