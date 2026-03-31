/**
 * Normalizes a file path from the API (which might have backslashes and redundant prefixes)
 * into a valid URL for the browser.
 */
export function getFileUrl(path: string | undefined | null): string {
    if (!path) return '';
    
    // 1. Normalize backslashes to forward slashes
    let normalized = path.replace(/\\/g, '/');
    
    // 2. Remove redundant /uploads prefix if it's already there
    // The backend seems to return paths starting with /uploads/ or \uploads\
    if (normalized.startsWith('/uploads/')) {
        normalized = normalized.substring(8); // remove /uploads
    } else if (normalized.startsWith('uploads/')) {
        normalized = normalized.substring(7); // remove uploads
    }
    
    // 3. Ensure it starts with a slash
    if (!normalized.startsWith('/')) {
        normalized = '/' + normalized;
    }
    
    // 4. Prepend the actual base uploads path
    return `/uploads${normalized}`;
}
