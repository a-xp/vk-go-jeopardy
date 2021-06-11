import config from './config';
import {fetch} from 'whatwg-fetch';

let vkParams;
try {
    vkParams = window.location.search.substring(1) + (window.location.hash ? '&game=' + window.location.hash.substring(1) : '');
} catch (e) {
    console.error(e);
}

function sendError(msg) {
    return fetch(config.api + '/log-error?msg=' + encodeURIComponent(msg));
}

async function makeRequest(method, url, params = {}) {

    if (params.query) {
        url = url + '?' + Object.keys(params.query).map(k => `${encodeURIComponent(k)}=${encodeURIComponent(params.query[k])}`).join('&');
    }

    const result = await fetch(config.api + url, {
        method,
        mode: 'cors',
        cache: 'no-cache',
        credentials: 'same-origin',
        headers: {
            'Content-Type': 'application/json',
            'Accept': 'application/json',
            'X-VK-PARAMS': vkParams
        },
        body: params.json ? JSON.stringify(params.json) : undefined
    })

    if (result.ok) {
        if (result.status !== 200 || parseInt(result.headers.get('Content-Length')) === 0) return {};
        return result.json();
    } else {
        throw await result.json();
    }
}

function getRating(from = 0, query = '') {
    return makeRequest('GET', '/api/rating', {query: {from, query}})
}

function getProfile() {
    return makeRequest('GET', '/api/me')
}

function getGames() {
    return makeRequest('GET', '/api/admin/games')
}

function getGame(id) {
    return makeRequest('GET', '/api/admin/games/' + id)
}

function deleteGame(id) {
    return makeRequest('DELETE', '/api/admin/games/' + id)
}

function updateGame(game) {
    return makeRequest('POST', '/api/admin/games', {json: game})
}

function getAdmins() {
    return makeRequest('GET', '/api/admin/admins')
}

function deleteAdmin(id) {
    return makeRequest('DELETE', '/api/admin/admins/' + id)
}

function addAdmin(link) {
    return makeRequest('POST', '/api/admin/admins', {json: {link}})
}

function getGroups() {
    return makeRequest('GET', '/api/admin/groups')
}

function deleteGroup(id) {
    return makeRequest('DELETE', '/api/admin/groups/' + id)
}

function addGroup(group) {
    return makeRequest('POST', '/api/admin/groups', {json: group})
}

const api = {
    getRating, getProfile, getGames, deleteGame, updateGame, getGame,
    getAdmins, addAdmin, deleteAdmin, getGroups, deleteGroup, addGroup,
    sendError
}

export default api;