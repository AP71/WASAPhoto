import {createApp, reactive} from 'vue'
import App from './App.vue'
import router from './router'
import axios from './services/axios.js'
import profile from './services/profile'
import ErrorMsg from './components/ErrorMsg.vue'
import LoadingSpinner from './components/LoadingSpinner.vue'
import Photo from './components/Photo.vue'
import NavBar from './components/NavBar.vue'

import './assets/main.css'


const app = createApp(App)
app.config.globalProperties.$axios = axios;
app.config.globalProperties.$profile = profile;
app.component("ErrorMsg", ErrorMsg);
app.component("LoadingSpinner", LoadingSpinner);
app.component("Photo", Photo);
app.component("NavBar", NavBar);
app.use(router)
app.mount('#app')
