<template>
  <div class="index container" v-scroll="onLoad">
    <top-list :articles="$store.state.articles" />
    <p v-if="isLoading" class="load-tip">加载中...</p>
    <p v-if="noMore" class="load-tip">没有更多文章了</p>
  </div>
    <!-- <div>
        <div class="home-categoties-box">
            <a href="/" class="categoties-item" :class="{'categoties-select': !cate}">全部</a>
            <a v-for="cateItem in categories" class="categoties-item" :href="'/?cate=' + cateItem.id" :class="{'categoties-select': cateItem.id == cate}">{{cateItem.name}}</a>
        </div>
        <div class="home-articles-box">
            <div v-for="article in articles" class="articles-cell">
                <a :href="'/user/' + article.user.id" target="_blank" class="user-icon-box"><img :src="article.user.avatarURL" alt=""/></a>
                <span class="home-tip-container">
                    <Tooltip :content="`回复数${article.commentCount}　浏览数${article.browseCount}`" placement="bottom-start" class="home-tip-box">
                        <a :href="'/topic/' + article.id" target="_blank" class="no-underline">
                            <span class="articles-click-num">{{article.commentCount}}</span>
                            <span class="articles-num-split">/</span>
                            <span class="articles-res-num">{{article.browseCount}}</span>
                        </a>
                    </Tooltip>
                </span>
                <span class="articles-categoties" :class="article.isTop ? 'articles-categoties-top' : 'articles-categoties-common' ">{{article.isTop ? '置顶' : article.categories[0].name}}</span>
                <a :href="'/topic/' + article.id" target="_blank" class="home-articles-title" :title="article.name">{{article.name | entity2HTML}}</a>
                <p class="articles-res-time">{{article.createdAt | getReplyTime}}</p>
                <a v-if="article.lastUser && article.lastUser.id" :href="'/user/' + article.lastUser.id" target="_blank" class="user-small-icon-box"><img :src="article.lastUser.avatarURL" alt=""/></a>
            </div>

            <div v-if="articles.length > 0" style="text-align: center;">
                <span v-if="totalVisible" class="ivu-page-total" style="margin-top: 10px;vertical-align: top;">共 {{totalCount}} 条</span>
                <Page class="common-page" :class="{'common-page-inline': totalVisible}"
                    :current="pageNo"
                    :page-size="pageSize"
                    :total="totalCount"
                    @on-change="onPageChange"
                    :show-elevator="true"/>
            </div>
        </div>
        <baidu-banner900x110 />
    </div> -->
</template>

<script>
    import request from '~/net/request'
    import dateTool from '~/utils/date'
    import htmlUtil from '~/utils/html'
    import baiduBanner900x110 from '~/components/ad/baidu/banner900x110'

    export default {
        data () {
            return {
                page: 1,
                isLoading: false,
                noMore: false
            }
        },
        asyncData (context) {
            context.store.commit('top10Visible', true)
            context.store.commit('friendLinkVisible', true)
            context.store.commit('statVisible', true)
            const query = context.query || {}
            return Promise.all([
                request.getCategories({
                    client: context.req
                }),
                request.getArticles({
                    client: context.req,
                    query: {
                        cateId: query.cate || '',
                        pageNo: query.pageNo || 1,
                        noContent: 'true'
                    }
                }),
                request.getTopList({
                    client: context.req
                })
            ]).then(data => {
                let user = context.user
                let cate = query.cate || ''
                let categories = data[0].data.categories || []
                let articles = data[1].data.articles || []
                let pageNo = data[1].data.pageNo
                let totalCount = data[1].data.totalCount
                let pageSize = data[1].data.pageSize
                let topList = (data[2] && data[2].data.articles) || []
                articles.map(items => {
                    items.isTop = false
                })
                if (!query.pageNo || parseInt(query.pageNo) < 2) {
                    topList.map(items => {
                        items.isTop = true
                    })
                    articles = topList.concat(articles)
                }
                return {
                    totalVisible: process.env.NODE_ENV !== 'production',
                    categories: categories,
                    articles: articles,
                    totalCount: totalCount,
                    pageNo: pageNo,
                    pageSize: pageSize,
                    user: user,
                    cate: cate
                }
            }).catch(err => {
                console.log(err.message)
                context.error({ message: 'Not Found', statusCode: 404 })
            })
        },
        head () {
            return {
                title: '首页'
            }
        },
        filters: {
            getReplyTime: dateTool.getReplyTime,
            entity2HTML: htmlUtil.entity2HTML
        },
        methods: {
            onPageChange (value) {
                window.location.href = `/?cate=${this.cate}&pageNo=${value}`
            },
            async onLoad () {
                // 没有正在加载中
                if (!this.isLoading) {
                    if (this.$store.state.articles.length < this.$store.state.total) {
                        this.isLoading = true
                        this.page++
                        await this.$store.dispatch('ARTICLES', this.page)
                        this.isLoading = false
                    } else {
                        this.noMore = true
                        return false
                    }
                } else {
                    return false
                }
            }
            // onLoad () {
            //     // 没有正在加载中
            //     if (!this.isLoading) {
            //         if (this.$store.state.articles.length < this.$store.state.total) {
            //             this.isLoading = true
            //             this.page++
            //             // await this.$store.dispatch('ARTICLES', this.page)
            //             request.getArticles({
            //                 // client: context.req,
            //                 query: {
            //                     // cateId: query.cate || '',
            //                     pageNo: this.page || 1
            //                     // noContent: 'true'
            //                 }
            //             }).then(res => {
            //                 // this.isLoading = false
            //                 // if (res.errNo === ErrorCode.SUCCESS) {
            //                 //     window.location.href = this.redirectURL
            //                 // } else if (res.errNo === ErrorCode.IN_ACTIVE) {
            //                 //     window.location.href = '/verify/mail?e=' + encodeURIComponent(res.data.email)
            //                 // } else {
            //                 //     // 没有配置luosimaoSiteKey的话，就没有验证码功能
            //                 //     this.luosimaoSiteKey && window.LUOCAPTCHA.reset()
            //                 //     this.$Message.error({
            //                 //         duration: config.messageDuration,
            //                 //         closable: true,
            //                 //         content: res.msg
            //                 //     })
            //                 // }
            //                 this.isLoading = false
            //             }).catch(err => {
            //                 this.isLoading = false
            //                 console.log(err.message)
            //                 // context.error({ message: 'Not Found', statusCode: 404 })
            //             })
            //             this.isLoading = false
            //         } else {
            //             this.noMore = true
            //             return false
            //         }
            //     } else {
            //         return false
            //     }
            // }
        },
        mounted () {
            console.log('golang123')
        },
        components: {
            'baidu-banner900x110': baiduBanner900x110
        }
    }
</script>

<style>
    @import '../assets/styles/home.css'
</style>
