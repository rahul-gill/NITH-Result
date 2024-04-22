<script>
    import { onMount } from 'svelte';
  
    let isDark = false;
  
    const toggleTheme = () => {
      isDark = !isDark;
      document.body.classList.toggle('dark', isDark);
      localStorage.setItem('theme', isDark ? 'dark' : 'light');
    };

    const setThemeFromPreference = () => {
        const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
        const savedTheme = localStorage.getItem('theme');
        if (savedTheme === 'dark' || (savedTheme === null && prefersDark)) {
            toggleTheme();
        }
    };
  
    onMount(() => {
        // Check for saved theme preference in localStorage
        setThemeFromPreference();

        // Listen for changes to user's preferred color scheme
        window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', setThemeFromPreference);

    });
</script>

<div class="theme-switch {isDark ? 'dark' : ''}">
    <input type="checkbox" class="checkbox" id="checkbox" on:change={toggleTheme} bind:checked={isDark}/>
    <label for="checkbox" class="label">
      <svg class="moon" width="24" height="24" stroke-width="1.5" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
        <path d="M3 11.5066C3 16.7497 7.25034 21 12.4934 21C16.2209 21 19.4466 18.8518 21 15.7259C12.4934 15.7259 8.27411 11.5066 8.27411 3C5.14821 4.55344 3 7.77915 3 11.5066Z" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"/>
      </svg>
      <svg class="sun" width="24" height="24" stroke-width="1.5" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
        <path d="M12 18C15.3137 18 18 15.3137 18 12C18 8.68629 15.3137 6 12 6C8.68629 6 6 8.68629 6 12C6 15.3137 8.68629 18 12 18Z" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"/>
        <path d="M22 12L23 12" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"/>
        <path d="M12 2V1" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"/>
        <path d="M12 23V22" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"/>
        <path d="M20 20L19 19" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"/>
        <path d="M20 4L19 5" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"/>
        <path d="M4 20L5 19" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"/>
        <path d="M4 4L5 5" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"/>
        <path d="M1 12L2 12" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"/>
      </svg>
      <div class="ball"></div>
    </label>
</div>

<style>
*, ::after, ::before {
  box-sizing: border-box;
}
.theme-switch {
  display: flex;
  justify-content: center;
  align-items: center;
  

    --switch-shadow: 0px 0px 10px 3px rgba(0, 0, 0, 0.1) inset;
    --transition: all 0.3s cubic-bezier(0.76, 0, 0.24, 1);
}

.theme-switch.dark {
    --switch-shadow: 0px 0px 10px 3px rgba(0, 0, 0, 0.5) inset
}

.theme-switch .checkbox {
  opacity: 0;
  position: absolute;
}


.theme-switch .label {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px;
  border: 1px solid var(--outline);
  border-radius: 50px;
  position: relative;
  height: 40px;
  width: 80px;
  cursor: pointer;
  box-shadow: var(--switch-shadow);
  transition: var(--transition);
}


.theme-switch .ball {
  transition: var(--transition);
  background-color: var(--on-background);
  position: absolute;
  border-radius: 50%;
  top: 5px;
  left: 5px;
  height: 30px;
  width: 30px;
 }

.theme-switch .moon {
  color: #f1c40f;
  transform-origin: center center;
  transition: all 0.5s cubic-bezier(0.76, 0, 0.24, 1);
  transform: rotate(0);
}

.theme-switch .sun {
  color: #ff6b00;
  transform-origin: center center;
  transition: all 0.5s cubic-bezier(0.76, 0, 0.24, 1);
  transform: rotate(0);
}

.theme-switch.dark .ball {
    transform: translatex(40px);
  }

</style>