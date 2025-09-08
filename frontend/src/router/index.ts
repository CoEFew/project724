import { createRouter, createWebHistory } from 'vue-router';
import Home from '../pages/Home.vue';
import DocumentsPage from '../pages/DocumentsPage.vue';
import PicturesPage from '../pages/PicturesPage.vue';
import MusicPage from '../pages/MusicPage.vue';
import VideosPage from '../pages/VideosPage.vue';

const routes = [
  { path: '/', component: Home },
  { path: '/documents', name: 'DocumentsPage', component: DocumentsPage },
  { path: '/pictures', name: 'PicturesPage', component: PicturesPage },
  { path: '/music', name: 'MusicPage', component: MusicPage },
  { path: '/videos', name: 'VideosPage', component: VideosPage },
];

export default createRouter({
  history: createWebHistory(),
  routes,
});
