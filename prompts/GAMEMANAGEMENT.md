# Game Management Logic To Ensure User Engagement

First, see glossary in ./RULES.md to understand the glossary terms.

Let's list out every action involving movement of cards:
- Dealing into the hand
- Playing a card (moving from one player's hand to In-Game pile)
- CUT (moving all the cards from the In-Game pile to the player who previous had the highest card in the In-Game suit before another player made the CUT)
- Discard (moving all the cards form the In-Game pile to a separate pile of face up cards away from the players or the In-Game pile)

The visualization of these actions must be handled in a way that engages the user.  Performing these actions quickly will result in confusion and must not happen.  Let's break down the visualization of how the actions should play out:

Deal - Start with the 52 cards, and set a 0.25s timer between dealing each card.  Deal selects a random player to start and then goes clockwise from there.  After a first card is dealt, pause for 0.25s before dealing the second card.  After a second card is dealt, pause for 0.25s before dealing the third card. And so on.  As and when a card is dealt, arrange it in the players' hand in the preferred sort order described in the rules.

Playing a card - When a player (human or bot) gets their turn, a card is selected and played.  This card needs to move from the hand to the In-Game pile.  Animate the move so that the card is shown slowly moving from the player or opponent hand to the pile and adjust the CSS gradually to go from player or opponent card style to the In-Game card style.  Once the animation is complete, before moving the turn to the next clockwise player, pause the game play for 3 seconds.

CUT - When a CUT is performed by a player, after the player's out-of-suit is visible on the In-Game pile, add a non-modal creative emphasizing indication of CUT. Then, before taking any action per the rules, pause the gameplay for 3 seconds.  Let everyone look at the CUT condition.  Add a note to the Session Updates as to who performed the CUT, and who had played the highcard before that got all the cards. After 3 seconds, animate the movement of the cards so that there's 0.25s pause after each card being moved to the player's hand.

Discard - When all the players have played a card to the In-Game pile and they're all the same suit, after the last player has played their hand and the card is visible in the In-Game pile, pause the gameplay for 3 seconds. No cards should be moved during these 3 seconds. Next, the cards need to move from the In-Game pile to the discard pile.  After 3 seconds, animate the movement of the cards so that there's 0.25s pause after each card being moved to the player's hand.



