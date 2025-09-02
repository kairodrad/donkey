# Comprehensive Rules

The donkey game ends with one loser out of 2-8 players.  The game has many rounds within a game.

# Glossary
Deal - when the 52 cards are shuffled and distributed among the players in a round-robin manner.
In-Game card pile - A set of cards on the board that does't belong to any player as they've just moved it from their hand to the board.
Play a card - when a user (human or bot) moves a card from their hand to the In-Game pile.
Ace of Spades - Anyone with Ace of Spades gets the first turn automatically, and must play that card onto the In-Game pile.
CUT - when an out-of-suit card is played by a player because they have no cards of the current suit on the In-Game pile.
Discard - When all players have played a card to the In-Game pile and they all match the suit, those cards are no longer in the game and are moved to a separate pile away from players and In-Game area.

# Round
- We start with 2-8 players having a distribution of all 52 cards.
- The objective is to get rid of all your cards within the rules of the game. When you get rid of all your cards, you're done for the round and are one of the winners.  There is only one loser per round who hasn't been able to get rid of all their cards. The loser accumulates a letter on every loss, starting with 'D' then 'O' then 'N' then 'K' then 'E' then 'Y'. When they accumulate all 6 letters over many rounds, we have identified the loser and the game ends.

## Start
- In order to start the game, we need 2-8 players.  The players may be a combination of real players and bot players.  Bots should not have any more information about the game state other than 1. knowing the game rules and their own cards, and 2. will play the game like an average player. Bots are assigned random one word names that have some humor in them, and the visualization will ensure to indicate that the player at the table is a bot player vs human.
- The first player to the URL starts a game.  This creates a unique game identifier. The player can copy a URL having this unique key and share it with other users.  The player may also bump up the bot count.  When the total players (human + bot) exceeds 2, the game may be started by the player by clicking a button.  Further, when the total exceeds 8, the game starts automatically.  As for the Game URL, the mode of sharing is not in scope of the game. Let's say it is shared via email, text, airdrop, whatever.
- While the status of every humar player connection is important and tracked, the first player keeps the whole game alive.  When the first player loses connection, the game is paused for everyone and a modal is put up for all users indicating the same.  When other human players lose connection, the game basically halts when it is that player's turn, but otherwise, there's no pausing with a modal.
- The game has a session updates chat window to keep in touch through messages as well as to look at the statuses of players.  You will see "Player X joined the game", "Player Y disconnected", "Player Z reconnected".
- Disconnected players, including the first player, can resume an incomplete game as long as they use the URL to connect to the game.  If the game ID in the URL is an active game, it'll make the player join the game.  Otherwise, the player will get a grid of active / incomplete games they are or were ever a part of, to rejoin.  The grid will not present game IDs but names of players plus bot count along with the last update time from the game in local MM/DD/YY HH:MI format.
- The first player can also abandon the game, which means that the game abruptly ends for everyone, and cannot be resumed by anyone.  It will never appear in the grid.
- Going back to sharing of game URL, the opponent player(s) who receive this unique URL will join the game directly, which increases the player count.  A 9th player joining the game will get an error as the 8th player joining would have started the game anyway.  If the first player creates a game and adds 4 bot players, then only 3 other players may join, and the 4th will get an error, but will be presented with the home page where they can create a new game and become the first player for that game.

## Deal
- The players are arranged in a circular sequence.  Against their position, you will eventually see their cards dealt for them, but also see their overall game status ("" or "D" or "DO" or "DONKE") to know how far away they are from losing. We will also use a red or green dot to indicate connection to the game.
- The UI should use a similar circular sequence to visualize the players.  For each player, the current player (themselves) should be at the head (say bottom center), and the rest in a circle around them.  The diagonally opposite person is visually above them in 2-D.
- The 52 cards are dealt in a clockwise sequence starting from a random player.
- Once the deal is done, only you can see your own cards.  In the visualization, you can see the backs of others' cards so that you can be aware of things like "How many cards do they have?". The cards may be arranged in a fan-hold.
- The cards are automatically arranged in a sorter order: suit-wise (Diamond then Clubs then Heart then Spades), and within the suit, number-wise in an increasing (2,3,4,5,6,7,8,9,10,J,Q,K,A).
- In terms of visualization, your own cards are paramount and are key to effectively playing the game.  Therefore, the visualization for your cards must be clear, and each card must be interactable in your view port.  It is only sufficient to see your opponents' cards.

## Gameplay
- The objective of the game is to get rid of all your cards per the rules.
- Whoever has the Ace of Spades starts the turn, and everyone goes clockwise from there.
- We will need a visualization of In-Game cards at the center of the table, which is a series of cards that each player has played during their turn.  These cards are always visible to all players (face up).  Since each player can only play one card before it is cleared, the In Game visualization can render the played card closer to the player who played it, so that it is clear to all users who has played what card.
- Opponents MUST match the suit of the card played by the first player.
- If the opponent DOES NOT have a card of the same suit, then they will perform a "CUT" of the suit by playing any preferred card from their hand.  The game engine must enforce that the CUT is legal.
- Once CUT, the turns end.  From the original suit of cards in the In-Game pile, We figure out who played the highest value card in that suit (From highest to lowest order: A, K, Q, J, 10, 9, 8, 7, 6, 5, 4, 3, 2). All the cards in the In-Game pile including the card used to CUT, goes to the hand of the player with the highest card.  For example: If Player 1 played 4 of Spades, and Player 2 played 8 of Spades and Player 3 cuts using Ace of Diamonds, then Player 2 will now own all 3 cards in their hand after the set of turn.
- After a CUT, the person who made the CUT will get the turn to play.
- It is also possible that all players have a card of the current suit in play.  In that case, the turns end when the last player in the sequence plays a card of the suit.  The player who went first will not go again.  When the last player in the sequence plays a card, the entire set of In-Game cards goes to a discard pile that will not feature in that round again.  If the round ends in a discard, the In-Game cards are reset to empty, and the person who played the highest value card starts the next set of turns.
- As the round iterates, some players will see their card count in hand grow, and for others, it will fall.  Anyone that exhausts all their cards is done for that round.  Note that it is possible that in one of the set of turns, you play your last card which is a high card, and then you get CUT later, which brings the whole In-Game pile into your hands.  So, you know you're surely done, when that set of turns end with a discard or with someone else playing a higher card than you.
- Finishing a round early doesn't mean anything. You're not a winner. It just means that you're not going to lose that round.
- The game then resumes with the remaining players and completely discounts the finished player.  The visualization will somehow indicate that a player is done for that round.
- The turn management system should be aware of finished players, and will assign turns leaving that player out.
- One by one, the players will start exhausting their cards.
- The player who is left with the final card, loses the round.  Their letter assignment goes up by 1 (D -> DO, DO -> DON, DONK -> DONKE, DONKE -> DONKEY). If the player ends up at DONKEY, the game completes.  This game, like abandoned games, cannot be resumed.  You will need to start a new game, which resets letters, etc., and starts everything over with the 1st Round.
- If there's no DONKEY, then we move on to the next round.  All players are in effect again, and we begin dealing cards again, and the same rules apply as the earlier round.  We repeat until we have a DONKEY.

## Strategy
- Given the rules, you can start thinking of strategies to win.
- High cards are risky, especially Aces.  On a CUT, you will get the entire pile.
- The risk of CUT is higher when the number of players is high.
- The risk of CUT is low if that suit hasn't had any cards go into discard pile already.
- As a player, you're incentivized to play high cards early on, as you know that holding on to them might be disadvantageous as the game goes on.
- This needs to be balanced with playing a high card and getting CUT early on, which can result in you collecting more high cards from other players.
- Players will benefit from remembering what they opponent collected when someone performed a CUT.  For eg, if a player who goes after you doesn't have Spades and you play Spades, they will CUT you and you pick up their card.  You will need to remember to not play Spades again when they're still in the game.  However, if the In-Game cards has a spade and the CUT cards went to that player, then you know that they now have a spade and you're safe to play a spade again. In this example, you are wiser to play a lower value Spade than what you know the opponent to have, so that you won't get cut by anyone else later on.

## Visualizations
- Let's not worry about animating the deal of the turns.  The card can teleport to the final locations.
- When the round starts, it can start with shuffled and dealt cards in everyone's hands.
- The person with Ace of Spades can get a dismissable notification that they have the starting card.
- You play a card:
  - On Desktop, but hovering over a selected card and clicking that card
  - On Mobile, tapping a card to select and tapping the same card to play. Note that tapping away from a selected card resets the selection.
- When you play a card, it teleports to the center of the table but visibly closer to your position in the center.  This card is now visible to all player.
- As the turns progress, the cards pile up at the center but visually subtly closer to the player that played that card.
- As the cards are played, highlight the highest card in the IN-Game pile, so that we know who gets the cards in case of a CUT.
- When CUT by a different suit, indicate a CUT through a small animation, and add to the Session Update log. The cards themselves can teleport to the hand of the person who had the highest card so far.
- When the turns don't end in a CUT and end in a discard, create a growing discard pile away from the table.  This discard pile is initially empty, but is later filled with cards stacked on top of each other.  Doesn't matter if they are face up or down.  It just remains there for the rest of the round.

## Bots
Let's create 3 types of bots that the first player can choose from as they are adding bots.
- Easy: These bots know the rules but sometimes don't use the best strategy to win.  They are sometimes ideal candidates to end up being a DONKEY.
- Medium: These bots know the rules very well and also adopt strategies to preserve their interests.  However, they do not try to remember others' cards to optimize their move.
- Difficult: These bots are better than medium because they do try to remember which player picked up which cards, and strategize their gameplay to improve their chances of winning (not being a DONKEY).  However, they do not have access to Game State data structure, just what is publicly visible in the In-Game pile.
