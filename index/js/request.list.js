$('#save-request').click(function () {
    let requestName = document.getElementById("input-request-name").value;
    if (requestName === "") {
        alert("request name is require")
    } else {
        let data = getReqResData();
        data.name = requestName
        insertRequest(data).then(success => {
            window.location.reload()
        }).catch(error => {
            alert(error);
        })
    }
});

$('#show-modal-save-request').click(function () {
    const reqData = getReqResData();
    const activeRequestName = getActiveRequestListName();
    if (activeRequestName === "") {
        console.log(activeRequestName);
        // generate name
        // name format will be method
        $('#input-request-name').val(`${reqData.selected_function}`)
        $('#saveRequest').modal('toggle');
    } else {
        reqData.id = activeRequestName.attributes["request-name"].value
        reqData.name = activeRequestName.innerText
        updateRequest(reqData).catch(error => {
            alert(error);
        })
    }
});
$('#show-modal-save-as-request').click(function () {
    const reqData = getReqResData();
    // generate name
    // name format will be method
    $('#input-request-name').val(`copy ${reqData.selected_function}`)
    $('#saveRequest').modal('toggle');
});

function getReqResData() {
    const serverTarget = document.getElementById("server-target").value;
    const selectService = document.getElementById("select-service").value;
    const selectFunction = document.getElementById("select-function").value;
    const responseHTML = document.getElementById("json-response").innerHTML;
    const schemaProtoHTML = document.getElementById("schema-proto").innerHTML;
    editor = ace.edit("editor");
    return {
        server_target: serverTarget,
        selected_service: selectService,
        selected_function: selectFunction,
        raw_request: editor.getValue(),
        response_html: responseHTML,
        schema_proto_html: schemaProtoHTML,
    }
}

function setReqResData(data) {
    $('#server-target').val(data.server_target);
    target = data.server_target;
    $("#select-service").html(new Option(data.selected_service, data.selected_service, true, true));
    $('#choose-service').show();
    $("#select-function").html(new Option(data.selected_function.substr(data.selected_service.length), data.selected_function, true, true));
    $('#choose-function').show();
    generate_editor(data.raw_request);
    $('#body-request').show();
    $('#schema-proto').html(data.schema_proto_html);
    $('#json-response').html(data.response_html);
    $('#response').show();
}

function resetReqResData() {
    target = "";
    $('#choose-service').hide();
    $('#choose-function').hide();
    $('#body-request').hide();
    $('#response').hide();
}

async function renderRequestList() {
    const ul = document.getElementById("request-list")
    ul.innerHTML = ""

    const nameList = await getAllRequestKey();
    nameList.forEach(function (item) {
        let node = document.createElement("li")
        node.classList.add("list-group-item", "request-list")
        node.setAttribute("request-name", item.id)
        // node.addEventListener("click", function (el) {
        //     console.log(el);
        //     updateRequestView(el.target.children[1])
        // });
        node.innerHTML = `
        <a title="Delete this request" class="delete-request" onclick="removeRequest('${item.id}')"><i class="fa fa-times"></i></a>
        <p class="one-long-line request" onclick="updateRequestView('${item.id}',this)">${item.name}</p>
        `
        ul.appendChild(node);
    })
}

function removeRequestSelectedClass() {
    const elems = document.querySelectorAll(".request-list");
    [].forEach.call(elems, function (el) {
        el.classList.remove("selected");
    });
}

function getActiveRequestListName() {
    const elems = document.querySelectorAll(".request-list");
    for (let i = 0; i < elems.length; i++) {
        const e = elems[i]
        if (e.classList.contains("selected")) {
            return e;
        }
    }
    return ""
}

function setServerTargetActive() {
    const elems = document.querySelectorAll('[for="server-target"]');
    [].forEach.call(elems, function (el) {
        el.classList.add("active");
    });
}

function updateRequestView(id,elm) {
    if (elm) {
        getRequest(id).then(data => {
            resetReqResData()
            setReqResData(data)
            removeRequestSelectedClass()
            elm.parentElement.classList.add('selected')
            setServerTargetActive();
        }).catch(error => {
            alert(error)
        })
    }
}

function removeRequest(id) {
    deleteRequest(id).then(() => {
        window.location.reload()
    }).catch((error) => {
        alert(error)
    })
}

function search(elm) {
    const li = document.querySelectorAll(".request-list")
    li.forEach(function (el) {
        if (el.getAttribute("request-name").toLowerCase().includes(elm.value.toLowerCase())) {
            el.style.display = ""
        } else {
            el.style.display = "none"
        }
    })
}

$(document).ready(function () {
    renderRequestList()
});