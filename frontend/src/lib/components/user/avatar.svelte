<script>
  let open = false;

  function toggleDropdown() {
    open = !open;
  }

  function closeDropdown(event) {
    if (!event.target.closest(".dropdown")) {
      open = false;
    }
  }

  window.addEventListener("click", closeDropdown);
</script>

<div class="dropdown">
  <button
    class="btn d-flex align-items-center text-decoration-none"
    on:click|stopPropagation={toggleDropdown}
  >
    <img
      src="placeholder.svg"
      alt="User Avatar"
      class="rounded-circle me-2"
      width="40"
      height="40"
    />
<i class="bi bi-caret-down-fill"></i>
  </button>

  {#if open}
    <ul class="dropdown-menu dropdown-menu-end show" style="position: absolute;">
      <li><a class="dropdown-item" href="/profile">Profile</a></li>
      <li><a class="dropdown-item" href="/settings">Settings</a></li>
      <li><hr class="dropdown-divider" /></li>
      <li><a class="dropdown-item text-danger" href="/logout">Logout</a></li>
    </ul>
  {/if}
</div>

<style>
  .dropdown {
    position: relative;
  }
  .dropdown-menu {
    display: block;
    top: 100%;
    margin-top: 0.25rem;
  }
</style>