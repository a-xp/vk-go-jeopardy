export function formatRusNumerics(n, forms) {
    const rem = n % 10;
    if (rem === 1) return forms[0];
    if (!rem || rem > 4) return forms[2];
    return forms[1];
}

export function formatPostId({groupId, postOwnerId, postId}) {
    return groupId ? `https://vk.com/club${groupId}?w=wall${postOwnerId}_${postId}` : '';
}

export function parsePostId(link) {
    const result = /vk\.com\/club(\d+)\?w=wall(-\d+)_(\d+)/.exec(link)
    return result ? {groupId: parseInt(result[1]), postOwnerId: parseInt(result[2]), postId: parseInt(result[3])} : null
}

export function hasText(str) {
    return typeof str === 'string' && /[0-9a-zа-я]+/.test(str)
}