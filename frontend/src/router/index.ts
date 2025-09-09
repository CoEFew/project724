import { createRouter, createWebHistory } from 'vue-router';
import Home from '../pages/Home.vue';
import DocumentsPage from '../pages/DocumentsPage.vue';
import CatText from '../pages/CatText.vue';
import MusicPage from '../pages/MusicPage.vue';
import VideosPage from '../pages/VideosPage.vue';

const routes = [
  { path: '/', component: Home },
  { path: '/DogPuzzle', name: 'DocumentsPage', component: DocumentsPage },
  { path: '/CatText', name: 'CatText', component: CatText },
  { path: '/music', name: 'MusicPage', component: MusicPage },
  { path: '/videos', name: 'VideosPage', component: VideosPage },
];

export default createRouter({
  history: createWebHistory(),
  routes,
});
