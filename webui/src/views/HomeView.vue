<script setup>
import Photo from '../components/Photo.vue'
import LoadingSpinner from '../components/LoadingSpinner.vue'
import ErrorMsg from '../components/ErrorMsg.vue'
</script>

<script>
export default {
	components: {
		Photo,
		LoadingSpinner,
		ErrorMsg,
	},
	data: function() {
		return {
			errormsg: null,
			loading: false,
			feed: [],
			nextPageId: 0,
		}
	},
	methods: {
		async getFeed() {
			this.loading = true;
			this.errormsg = null;
			try{
				const nextPage = this.nextPageId;
				let response = nextPage == 0?  await this.$axios.get("/feed/") : await this.$axios.get(`/feed/?pageId=${this.nextPageId}`) ;
				if (!(response.status === 204)) {
					nextPage == 0 ? this.feed = response.data.posts : this.feed.push(...response.data.posts);
					this.nextPageId = response.data.nextFeedPageId;
				}
			} catch(e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
  		},
	},
	mounted() {
		if (this.$profile.identifier==null){
			this.$router.push({path: '/'});
		} else {
			this.getFeed();
		}
	}
}
</script>

<template>
	<div id="home" class="d-flex min-vh-100 w-100 justify-content-center align-items-center" style="background-color: #383838">
		<div class="d-flex flex-column justify-content-center align-items-center min-vh-100 w-75" style="background-color: #2e2e2e;">		
			<div style="height: 60px"/>
			<LoadingSpinner :loading="this.loading"/>
			<ErrorMsg :msg="this.errormsg"/>
			<Photo v-for="post in this.feed" v-bind:key="post.photo" v-bind:post="post"/>
			<div v-if="this.nextPageId==0" class="d-flex flex-row justify-content-center align-items-center p-4">
				<div class="rounded rounded-5 fs-5 text-success py-2 px-5" style="background-color: #212121;">
					Non c'è più niente qui!! Vai a comprare un gelato...
				</div>
			</div>
			<div v-else class="p-4">
				<button type="button" class="btn btn-outline-success text-white fw-bolder rounded-pill fs-4" style="width: 150px" @click="this.getFeed">...</button>
			</div>
			<div style="height: 75px"/>
		</div>
	</div>
</template>

<style>
</style>
