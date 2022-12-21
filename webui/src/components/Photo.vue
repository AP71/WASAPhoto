<script setup>
import LoadingSpinner from '../components/LoadingSpinner.vue'
import ErrorMsg from '../components/ErrorMsg.vue'
</script>

<script>
export default {
	props: ['post'],
	components: {
		LoadingSpinner,
		ErrorMsg,
	},
    data: function() {
		return {
			errormsg: null,
			loading: null,
			blob: null,
			blobUrl: null,
			showCommentStatus: false,
			liked: false,
			comments: [],
			nextCommentPageId: 0,
			newComment: null,
		}
	},
	methods: {
		async getPhoto() {
			this.loading = true;
			this.errormsg = null;
			try{
				let response = await this.$axios.get(`/feed/${this.post.photo}/`, { responseType: 'blob'});
				this.blob = response.data
				this.blobUrl = URL.createObjectURL(this.blob);
				document.getElementById('post-image').src = this.blobUrl;
			} catch(e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
  		},
		async getStatus() {
			this.loading = true;
			this.errormsg = null;
			// status
			try{
				let response = await this.$axios.get(`/feed/${this.post.photo}/likes/${this.$profile.username}`);
				this.liked = true;
			} catch(e) {
				if (e.response.status == 404) {
					this.liked = false;
				} else {
					this.errormsg = e.toString();
				}
			}
			this.loading = false;
  		},
		async setLike() {
			this.loading = true;
			this.errormsg = null;
			try{
				let response = await this.$axios.put(`/feed/${this.post.photo}/likes/${this.$profile.username}`);
				this.liked = true;
				this.post.numberOfLikes += 1;
			} catch(e) {
				this.errormsg = e.toString();
				if (e.response.status == 409) {
					this.liked = true;
				}
			}
			this.loading = false;
  		},
		async removeLike() {
			this.loading = true;
			this.errormsg = null;
			try{
				let response = await this.$axios.delete(`/feed/${this.post.photo}/likes/${this.$profile.username}`);
				this.liked = false;
				this.post.numberOfLikes -= 1;
			} catch(e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
  		},
		showComment() {
			this.showCommentStatus = !this.showCommentStatus;
		},
		async deleteComment(commentId) {
            this.loading = true;
			this.errormsg = null;
			try{
				let response = await this.$axios.delete(`/feed/${this.post.photo}/comments/${this.$profile.username}?idCommento=${commentId}`);
				this.post.numberOfComments -= 1;
			} catch(e) {
				this.errormsg = e.toString();
			}
			this.getComments();
        },
		async sendComment(commentId) {
			if (this.newComment != null) {
				this.loading = true;
				this.errormsg = null;
				try{
					let response = await this.$axios.post(`/feed/${this.post.photo}/comments/${this.$profile.username}`, {text: this.newComment});
					this.post.numberOfComments += 1;
					this.newComment = null
				} catch(e) {
					this.errormsg = e.toString();
				}
				this.getComments();
			}
        },
		async getComments() {
			this.loading = true;
			this.errormsg = null;
			try{
				let response = await this.$axios.get(`/feed/${this.post.photo}/comments/`);
				this.comments = response.data.comments
				this.nextCommentPageId = response.data.nextCommentPageId;
			} catch(e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		}
	},
	mounted() {
		this.getPhoto();
		this.getStatus();
		this.getComments();
	},
}
</script>

<template>
	<LoadingSpinner :loading="this.loading"/>
	<ErrorMsg :msg="this.errormsg"/>
	<div class="card" style="background-color: #212121;">
        <div class="PhotoHeader">
			<div class="username fw-bold fs-3">
				{{ this.post.username }}
			</div>
			<div class="data fs-6"> 
				{{ this.post.data }}
			</div>
		</div>		
		<div class="d-flex justify-content-center align-items-center">
			<img id="post-image" class="image">
		</div>
		<div class="Buttons">
			<div class="d-flex justify-content-start align-items-center min-vw-25 max-vw-25">
				<div class="text-white fs-bold fs-4 pe-2">
					{{ this.post.numberOfLikes }}
				</div>
				<i v-if="this.liked" class="icon bi bi-heart-fill" @click="removeLike" style="color:red"></i>
				<i v-else class="icon bi bi-heartbreak" @click="setLike" style="color:red"></i>
			</div>
			<div class="d-flex justify-content-center align-items-center min-vw-25 max-vw-25">
				<div class="text-white fs-bold fs-4 pe-2">
					{{ this.post.numberOfComments }}
				</div>
				<i class="icon bi bi-chat-left-text" @click="showComment"></i>
			</div>
		</div>	
		<div class="commentAreaInput">
			<textarea v-model="this.newComment" class="textArea border border-2 border-success" placeholder="Commento" style="background-color: #212121;"></textarea>
			<i class="icon bi bi-send-fill text-success fs-4" @click="sendComment"></i>
		</div>
		<div v-if="this.showCommentStatus" class="commentArea">
			<div v-for="comment in this.comments" key="comment" class="comment pb-2">
				<div class="dettagli">
					<div class="fw-bolder fs-5 text-white pe-4">
						{{ comment.username }}
					</div>
					<div class="fw-normal fs-6 text-white">
						{{ comment.text }}
					</div>
				</div>
				<div class="dettagli">
					<div class="text-bold text-success fs-6 pe-4">
						{{ comment.data }}
					</div>
					<div v-if="comment.idUser==this.$profile.identifier || this.post.identifier==this.$axios.identifier" class="text-decoration-underline text-success fs-6" @click="deleteComment(comment.id)">
						Elimina
					</div>
				</div>
			</div>
		</div>
    </div>
</template>

<style scoped>
	.card {
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		padding-top: 2rem;
		width: 60%;
		border-radius: 1rem;
	}
	.PhotoHeader {
		display: flex;
		justify-content: space-between;
		align-items: center;
		width: 90%;
	}
	.username {
		display: flex;
		justify-content: start;
		align-items: center;
		width: 90%;
		padding-bottom: 1rem;
		color: white;
	}
	.image {
		object-fit: contain;
		width: 90%;
		padding-bottom: 0.5rem;
	}
	.data {
		display: flex;
		justify-content: right;
		align-items: center;
		width: 90%;
		padding-bottom: 1rem;
		color: white;
	}
	.Buttons {
		display: flex;
		justify-content: space-between;
		align-items: center;
		width: 90%;
	}

	.commentAreaInput {
		display: flex;
		justify-content: space-between;
		align-items: center;
		width: 90%;
		height: 10%;
		padding-bottom: 1rem;
	}

	.textArea {
		min-width: 92%;
		border-radius: 2rem;
		height: 2.25rem;
		color: white;
		text-indent: 10px;
	}

	.commentArea {
		display: flex;
		flex-direction: column;;
		justify-content: space-between;
		align-items: center;
		width: 90%;
		overflow: scroll;
		max-height: 150px;
		padding-bottom: 1rem;
	}

	.comment {
        display: flex;
        flex-direction: column;
        justify-content: center;
        align-items: center;
        width: 100%;
    }

    .dettagli {
        display: flex;
        flex-direction: row;
        justify-content: start;
        align-items: center;
        width: 100%;
    }

</style>

