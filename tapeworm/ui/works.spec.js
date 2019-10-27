// const fs = require('fs');

describe('it works?', () => {
    it('should be working', () => {
        expect(process.version).toBe('v13.0.1');
    })
})

describe('bundling', () => {
    it('should work', () => {
        let written;

        console.log = (m) => written = m;
        const bundle = require('com_github_jaxsax_projects/tapeworm/ui/bundle.js');
        expect(written).toEqual('Hello, Bob');
    })
})
