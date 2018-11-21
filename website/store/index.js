import Vue from 'vue'
import Vuex from 'vuex'
import getters from './getters'
// import mutations from './mutations'
import actions from './actions'

Vue.use(Vuex)

const store = () => new Vuex.Store({
    state: {
        siteConfig: {},
        baiduAdConfig: {},
        isAdminPage: false,
        user: null,
        messages: [],
        messageCount: 0,
        userCount: 0,
        topicCount: 0,
        replyCount: 0,
        top10Users: [],
        maxCommentArticles: [],
        maxBrowseArticles: [],
        topicAuthor: null,
        authorRecentArticles: [],
        votesMaxBrowse: [],
        votesMaxComment: [],

        userLoginVisible: true,
        publishTopicVisible: true,
        createBookVisible: false,
        top10Visible: false,
        maxBrowseVisible: true,
        maxCommentVisible: true,
        friendLinkVisible: false,
        statVisible: false,
        token: '',
        githubToken: '',
        // user: {},
        total: 0,
        articles: [],
        tags: [],
        tagArticles: [], // 标签下的文章列表
        searchArticles: [],
        archives: [],
        githubApi: {
            userInfo: 'https://api.github.com/user?access_token='
        }
    },
    getters,
    mutations: {
        siteConfig (state, siteConfig) {
            state.siteConfig = siteConfig
        },
        baiduAdConfig (state, baiduAdConfig) {
            state.baiduAdConfig = baiduAdConfig
        },
        isAdminPage (state, isAdminPage) {
            state.isAdminPage = isAdminPage
        },
        messages (state, messages) {
            state.messages = messages
        },
        messageCount (state, messageCount) {
            state.messageCount = messageCount
        },
        user (state, user) {
            state.user = user
        },
        avatarURL (state, url) {
            if (state.user) {
                state.user.avatarURL = url
            }
        },
        top10Users (state, top10Users) {
            state.top10Users = top10Users
        },
        maxCommentArticles (state, maxCommentArticles) {
            state.maxCommentArticles = maxCommentArticles
        },
        maxBrowseArticles (state, maxBrowseArticles) {
            state.maxBrowseArticles = maxBrowseArticles
        },
        topicAuthor (state, topicAuthor) {
            state.topicAuthor = topicAuthor
        },
        authorRecentArticles (state, authorRecentArticles) {
            state.authorRecentArticles = authorRecentArticles
        },
        userLoginVisible (state, userLoginVisible) {
            state.userLoginVisible = userLoginVisible
        },
        publishTopicVisible (state, publishTopicVisible) {
            state.publishTopicVisible = publishTopicVisible
        },
        createBookVisible (state, createBookVisible) {
            state.createBookVisible = createBookVisible
        },
        top10Visible (state, top10Visible) {
            state.top10Visible = top10Visible
        },
        maxBrowseVisible (state, maxBrowseVisible) {
            state.maxBrowseVisible = maxBrowseVisible
        },
        maxCommentVisible (state, maxCommentVisible) {
            state.maxCommentVisible = maxCommentVisible
        },
        friendLinkVisible (state, friendLinkVisible) {
            state.friendLinkVisible = friendLinkVisible
        },
        statVisible (state, statVisible) {
            state.statVisible = statVisible
        },
        userCount (state, userCount) {
            state.userCount = userCount
        },
        topicCount (state, topicCount) {
            state.topicCount = topicCount
        },
        replyCount (state, replyCount) {
            state.replyCount = replyCount
        },
        SET_TOKEN (state, token) {
            state.token = token
        },
        SET_GITHUB_TOKEN (state, token) {
            state.githubToken = token
        },
        SET_USER (state, data) {
            state.user = data.data
        },
        SET_ARTICLES (state, data) {
            state.articles = state.articles.concat(data.data)
            state.total = data.total
        },
        SET_TAGS (state, data) {
            state.tags = data.data
        },
        SET_TAG_ARTICLES (state, data) {
            state.tagArticles = data.data
        },
        SET_ARCHVES (state, data) {
            state.archives = data.data
        },
        SET_SEARCH (state, data) {
            state.searchArticles = data.data
        }
    },
    actions
})

export default store
