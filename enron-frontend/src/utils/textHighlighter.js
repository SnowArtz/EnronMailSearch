export const highlightText = (text, searchTerms) => {
    if (!text || !searchTerms) return text;

    let terms = [];

    if (searchTerms.query) {
        terms.push(...searchTerms.query.split(' '));
    }
    if (searchTerms.from) terms.push(searchTerms.from);
    if (searchTerms.to) terms.push(searchTerms.to);

    terms = terms.filter(term => term && term.trim() !== '');

    const escapeRegExp = (string) => {
        return string.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
    };

    if (terms.length > 0) {
        const pattern = new RegExp(terms.map(term => `(${escapeRegExp(term)})`).join('|'), 'gi');
        return text.replace(pattern, match => `<mark>${match}</mark>`);
    }

    return text;
}; 