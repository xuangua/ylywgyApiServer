import cookie from '~/utils/cookie'

export default ({ app: { router }, req, res }) => {
    router.afterEach((to, from) => {
        if (typeof window === 'undefined') {
            cookie.refreshTokenCookie(req, res)
        }
    })
    // console.log(req)
    // console.log(res)
    if (!process.server) {
        console.log('middleware from client side')
    } else {
        console.log('middleware from server side')
    }
}
