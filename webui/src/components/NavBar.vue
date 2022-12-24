<script>
export default {
	methods: {
		setStatus(i) {
			if (i==0) {
				this.$router.push({path: '/home'});
				document.getElementById("1").setAttribute('class', 'icon bi bi-house text-success');
				document.getElementById("2").setAttribute('class', 'icon bi bi-search text-white');
				document.getElementById("3").setAttribute('class', 'icon bi bi-person text-white');
			}
			else if (i==1) {
				this.$router.push({path: '/search'});
				document.getElementById("1").setAttribute('class', 'icon bi bi-house text-white');
				document.getElementById("2").setAttribute('class', 'icon bi bi-search text-success');
				document.getElementById("3").setAttribute('class', 'icon bi bi-person text-white');
			}
			else if (i==2) {
				this.$router.push({path: `/profiles/${this.$profile.username}`});
				document.getElementById("1").setAttribute('class', 'icon bi bi-house text-white');
				document.getElementById("2").setAttribute('class', 'icon bi bi-search text-white');
				document.getElementById("3").setAttribute('class', 'icon bi bi-person text-success');
			}
		}
	},
}
</script>

<template>
	<div class="navbar">
		<div class="element" @click="setStatus(0)">
            <i id="1" class="icon bi bi-house text-success"></i>
        </div>
        <div class="element" @click="setStatus(1)">
            <i id="2" class="icon bi bi-search text-white"></i>
        </div>
        <div class="element" @click="setStatus(2)">
            <i id="3" class="icon bi bi-person text-white"></i>
        </div>
	</div>
</template>

<style>
    .icon {
		color: white;
        font-size: 2rem;
    }
    .navbar {
		display: flex;
		flex-direction: row;
		justify-content: space-evenly;
		align-items: center;
		background-color: #212121;
		overflow: hidden;
		position: fixed;
		bottom: 0;
		left: 50%;
		transform: translateX(-50%);
		width: 75%;
		height: 75px;
	}
	.element {
		display: flex;
		justify-content: center;
		align-items: center;
		min-width: 25%;
	}
</style>
