module.exports = {
    head: {
        meta: [
            { charset: 'utf-8' },
            { name: 'viewport', content: 'width=device-width, initial-scale=1' }
        ],
        link: [
            { rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' },
            { rel: 'alternate', type: 'application/rss+xml', title: 'RSS 2.0', href: '/rss.xml' },
            { rel: 'stylesheet', type: 'text/css', href: '//at.alicdn.com/t/font_620064_hrymm1e94nnlv7vi.css' },
            { rel: 'stylesheet', href: '/styles/iview-2.11.0.css' }
        ]
    },
    css: [
        '~assets/styles/common.css','~assets/css/main.scss'
    ],
    loading: { color: '#80bd01' },
    performance: {
        prefetch: false
    },
    render: {
        resourceHints: false
    },
    build: {
        /*
         ** Run ESLINT on save
         */
        extend (config, ctx) {
            if (ctx.isClient) {
                config.module.rules.push({
                    enforce: 'pre',
                    test: /\.(js|vue)$/,
                    loader: 'eslint-loader',
                    exclude: /(node_modules)/
                })
            }
        },
        vendor: ['axios', 'iview']
    },
    modules: ['@nuxtjs/pwa', '@nuxtjs/axios'],
    plugins: [
        { src: '~plugins/iview.js', ssr: true },
        { src: '~plugins/bdStat.js', ssr: false },
        { src: '~plugins/analyze.js', ssr: false },
        // { src: '~plugins/adsense.js', ssr: false },
        { src: '~/plugins/components.js', ssr: true},
        { src: '~/plugins/filters.js', ssr: true},
        { src: '~plugins/refreshToken.js', ssr: true }
    ]
}
