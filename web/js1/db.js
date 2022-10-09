function getRequest(name) {
    return new Promise(async function (success, error) {
        const resp = await fetch(`/api/request/${name}`)
        const data = await resp.json()
        success(data.data)
    })
}

function getAllRequestKey() {
    return new Promise(async function (success, error) {
        const resp = await fetch('/api/requests')
        const data = await resp.json()
        success(data.data)
    })
}

function insertRequest(request) {
    return new Promise(async function (success, error) {
        const resp = await fetch('/api/request', {
            method: 'POST', headers: {
                'Content-Type': 'application/json',
            }, body: JSON.stringify(request),
        })

        success(resp.ok)
    })
}

function updateRequest(request) {
    return new Promise(async function (success, error) {
        const resp = await fetch(`/api/request/${request.id}`, {
            method: 'PUT', headers: {
                'Content-Type': 'application/json',
            }, body: JSON.stringify(request),
        })

        success(resp.ok)
    })
}

function deleteRequest(name) {
    return new Promise(async function (success, error) {
        const resp = await fetch(`/api/request/${name}`, {method: 'DELETE'})
        success(resp.ok)
    })
}