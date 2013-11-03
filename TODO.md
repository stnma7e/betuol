TODO
====

<u>_Short-term_</u>
----------

- Fix graphics organization. Make a system for switching between models and shaders without sending all uniform values regardless of the shader. Be able to determine the values that need to be sent before rendering (make shader/program class and this renders instead?).

- Character system
	- Interaction through discrete events: each character has a list of events invovled with and chooses one to interact with. AI chooses at random/decision tree for best outcome, and the players choose at will? Will this remove capability for advanced interactions later? I can always change later if it doesn't work out.

	- Type of interaction:
		- Attack
		- Run: removes interaction after certain distance.

- Quest sytem

- Convert character system to Erlang-style goroutine based actors. Each actor is responsible for listening to events that it is interested in/responding to character interactions. Since goroutines are supposed to be lightweight, and you can have multiple Instances (potentially across many computers) there shouldn't be too much overhead. This will help with the Quest/AI/Chat systems a lot.

- Fix 4x4 matrix inverse

<u>_Long-term_</u>
---------

- Collision/physics system.

#### Optimizations

- Shaders are expensive to change. Sort rendering so that the least amount of shaders/programs are changed.

<u>_Story_</u>
-------------

- Character attributes are leveled based on activity. If a character uses strength a lot, then the character will level his strength.

- Multiple areas for attacking: head, chest, arms, legs? The damage done will be greater based on the chance of hitting: chest -> legs -> arms -> head.

- You can tell the attributes of a player by size/external features.

- Players are in one of two factions: upper class, lower class (name later). Upper class lives in acropolis and has money/looks nice. Lower class lives outside city (in the village/slums). Each has it's own community and players.

- Collecting taxes from players. Tax collector comes around every so often and asks for taxes. If you don't pay, then the AI city officials will hunt you down.

- Patron-client system.

- Metallurgy: players make their own weapons/tools. Can mine ore and smelt it/make alloys and mix molten metals to make stronger ones.
   - Metals have attributes that are combined in certain ways when mixed.

- Open PVP: Players can attack anyone if the attacker is a lower level than what he is attacking. If not AI city guards etc. come attack the higher level guy.

- Guilds/clubs for classes.

- Three (or more) warring factions. When in a state of peace, trade is allowed and players can interact, when in a state of war, player can interact but only "underground", in secret, away from the NPCs of the factions. 1984-esque


