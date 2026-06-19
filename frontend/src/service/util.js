export function formatRusNumerics(n, forms) {
    const rem = n % 10;
    if (rem === 1) return forms[0];
    if (!rem || rem > 4) return forms[2];
    return forms[1];
}

export function formatPostId({postOwnerId, postId}) {
    return postOwnerId ? `https://vk.com/wall${postOwnerId}_${postId}` : '';
}

export function parsePostId(link) {
    let result = /vk\.com\/club\d+\?w=wall(-\d+)_(\d+)/.exec(link);
    if (result) {
        return {postOwnerId: parseInt(result[1]), postId: parseInt(result[2])};
    }
    result = /vk\.com\/wall(-\d+)_(\d+)/.exec(link);
    return result ? {postOwnerId: parseInt(result[1]), postId: parseInt(result[2])} : null
}

export function hasText(str) {
    return typeof str === 'string' && /[0-9a-zа-я]+/.test(str)
}